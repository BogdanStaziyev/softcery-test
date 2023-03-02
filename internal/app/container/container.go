package container

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/database"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/http/controllers"
	"github.com/BogdanStaziyev/softcery-test/internal/rabbit"
	"github.com/BogdanStaziyev/softcery-test/internal/service"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
)

// Container holds the values of controller and queue broker services
// which can be extended with middleware, etc.
type Container struct {
	Services
	Controllers
	Queue
}

// Services struct of all services
type Services struct {
	service.ImageService
}

// Controllers struct of all handlers
type Controllers struct {
	controllers.ImageHandler
}

// Queue struct of all queue now we use rabbitMQ
type Queue struct {
	rabbit.Rabbit
}

// New куегкт all the dependencies required for the application to work as described in the above structures
func New(conf config.Configuration) Container {
	//PostgreSQL session
	sess := getDbSess(conf)

	//Create a new connection to RabbitMQ
	rabbitMQ := rabbit.NewRabbit(conf.RabbitURL)

	//Create image repository
	imageRepo := database.NewImageRepo(sess)

	//Create image service
	imageService := service.NewImageService(conf.FileStorageLocation, imageRepo, rabbitMQ)

	//Create image handler
	imageHandler := controllers.NewImageHandler(imageService)

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

// getDbSess create session with PostgreSQL.
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
