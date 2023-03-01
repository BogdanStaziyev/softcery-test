package container

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/app"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/database"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/handlers"
	"github.com/BogdanStaziyev/softcery-test/internal/rabbit"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
)

type Container struct {
	Services
	Controllers
	Queue
}

type Services struct {
	app.ImageService
}

type Controllers struct {
	handlers.ImageHandler
}

type Queue struct {
	rabbit.Rabbit
}

func New(conf config.Configuration) Container {
	//PostgreSQL session
	sess := getDbSess(conf)

	//Create a new connection to RabbitMQ
	rabbitMQ := rabbit.NewRabbit(conf.RabbitURL)

	//Create queue
	err := rabbitMQ.CreateQueue()
	if err != nil {
		log.Fatalf("RabbitMQ create queue error: %q\n", err)
	}

	//Create a consumer that continuously reads messages containing image path.
	//Forwards the path to create different versions of the photo.
	go func() {
		err = rabbitMQ.Consumer()
		if err != nil {
			log.Fatalf("RabbitMQ consumer error: %q\n", err)
		}
	}()

	//Create image repository
	imageRepo := database.NewImageRepo(sess)

	//Create image service
	imageService := app.NewImageService(conf.FileStorageLocation, imageRepo, rabbitMQ)

	//Create image handler
	imageHandler := handlers.NewImageHandler(imageService)

	return Container{
		Services: Services{
			imageService,
		},
		Controllers: Controllers{
			imageHandler,
		},
		Queue: Queue{
			rabbitMQ,
		},
	}
}

// Create session with PostgreSQL.
func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
