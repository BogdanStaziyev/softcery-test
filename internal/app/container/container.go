package container

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/controller/rabbit"
	"github.com/BogdanStaziyev/softcery-test/internal/usecase/database"
	"github.com/BogdanStaziyev/softcery-test/internal/usecase/service"
	session "github.com/BogdanStaziyev/softcery-test/pkg/database"
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
)

// Container holds the values of controller and queue broker services
// which can be extended with middleware, etc.
type Container struct {
	Services
	Queue
}

// Services struct of all services
type Services struct {
	service.ImageService
}

// Queue struct of all queue now we use rabbitMQ
type Queue struct {
	rabbit.Rabbit
}

// New куегкт all the dependencies required for the application to work as described in the above structures
func New(conf config.Configuration, l logger.Interface) Container {
	//PostgreSQL session
	sess := session.NewDbSess(&session.Config{
		DatabaseUser:     conf.DatabaseUser,
		DatabaseName:     conf.DatabaseName,
		DatabaseHost:     conf.DatabaseHost,
		DatabasePassword: conf.DatabasePassword,
	})

	//Create a new connection to RabbitMQ
	rabbitMQ := rabbit.NewRabbit(conf.RabbitURL, l)

	//Create image repository
	imageRepo := database.NewImageRepo(sess)

	//Create image service
	imageService := service.NewImageService(conf.FileStorageLocation, imageRepo, rabbitMQ)

	return Container{
		Services: Services{
			imageService,
		},
		Queue: Queue{
			rabbitMQ,
		},
	}
}
