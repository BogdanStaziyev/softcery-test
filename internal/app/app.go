package app

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/app/container"
	"github.com/BogdanStaziyev/softcery-test/internal/controller/http/v1"
	"github.com/BogdanStaziyev/softcery-test/pkg/httpserver"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(conf config.Configuration) {
	//make migration
	err := Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %s\n", err)
	}

	//initialize storage location
	_, err = os.Stat(conf.FileStorageLocation)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(conf.FileStorageLocation, os.ModePerm)
		if err != nil {
			log.Fatalf("storage folder can not be created, %s", err)
		}
	} else if err != nil {
		log.Fatalf("storage folder is not available %s", err)
	}

	//initialize container.go with delete services and db
	cont := container.New(conf)

	//Create queue
	err = cont.Rabbit.CreateQueue()
	if err != nil {
		log.Fatalf("RabbitMQ create queue error: %q\n", err)
	}

	//Create a consumer that continuously reads messages containing image path.
	//Forwards the path to create different versions of the photo.
	go func() {
		err = cont.Rabbit.Consumer()
		if err != nil {
			log.Fatalf("RabbitMQ consumer error: %q\n", err)
		}
	}()

	// HTTP Server
	handler := echo.New()
	v1.EchoRouter(handler, cont)
	httpServer := httpserver.New(handler, conf.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Print(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Print(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
