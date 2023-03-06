// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/BogdanStaziyev/softcery-test/internal/app/container"
	"github.com/BogdanStaziyev/softcery-test/internal/controller/http/v1"

	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/pkg/httpserver"
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
)

func Run(conf config.Configuration) {
	l := logger.New("debug")

	//make migration
	err := Migrate(conf)
	if err != nil {
		l.Fatal("Unable to apply migrations: ", "err", err)
	}

	//initialize storage location
	_, err = os.Stat(conf.FileStorageLocation)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(conf.FileStorageLocation, os.ModePerm)
		if err != nil {
			l.Fatal("storage folder can not be created", "err", err)
		}
	} else if err != nil {
		l.Fatal("storage folder is not available", "err", err)
	}

	//initialize container.go with delete services and db
	cont := container.New(conf, l)

	//Create queue
	err = cont.Queue.CreateQueue()
	if err != nil {
		l.Fatal("Queue create queue error: ", "err", err)
	}

	//Create a consumer that continuously reads messages containing image path.
	//Forwards the path to create different versions of the photo.
	go func() {
		err = cont.Queue.Consumer()
		if err != nil {
			l.Fatal("Queue consumer error: ", "err", err)
		}
	}()

	// HTTP Server
	handler := echo.New()
	v1.EchoRouter(handler, cont.Services, l)
	httpServer := httpserver.New(handler, conf.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Error("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: ", "err", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown: ", "err", err)
	}
}
