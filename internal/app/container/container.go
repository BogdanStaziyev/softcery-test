package container

import (
	// config
	"github.com/BogdanStaziyev/softcery-test/config"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/queue"
	"github.com/BogdanStaziyev/softcery-test/internal/usecase"
	"github.com/BogdanStaziyev/softcery-test/internal/usecase/database"

	// external
	session "github.com/BogdanStaziyev/softcery-test/pkg/database"
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
)

// Container holds the values of queue, broker, services
type Container struct {
	Services
	Queue
}

// Services struct of all services
type Services struct {
	usecase.ImageService
}

// Queue struct of all queue now we use rabbitMQ
type Queue struct {
	usecase.Queue
}

// New return all the dependencies required for the application to work as described in the above structures
func New(conf config.Configuration, l logger.Interface) Container {
	// postgreSQL session
	sess := session.NewDbSess(&session.Config{
		DatabaseUser:     conf.DatabaseUser,
		DatabaseName:     conf.DatabaseName,
		DatabaseHost:     conf.DatabaseHost,
		DatabasePassword: conf.DatabasePassword,
	})

	//Create a new connection to RabbitMQ
	rabbitMQ := queue.NewRabbit(conf.RabbitURL, l)

	//Create image repository
	imageRepo := database.NewImageRepo(sess)

	//Create image service
	imageService := usecase.NewImageService(conf.FileStorageLocation, imageRepo, rabbitMQ)

	return Container{
		Services: Services{
			imageService,
		},
		Queue: Queue{
			rabbitMQ,
		},
	}
}
