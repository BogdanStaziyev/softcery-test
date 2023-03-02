package app

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/app/container"
	myHttp "github.com/BogdanStaziyev/softcery-test/internal/infra/transport"
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

	//initialize container.go with handlers services and db
	cont := container.New(conf)

	// HTTP Server
	handler := echo.New()
	myHttp.EchoRouter(handler, cont)
	httpServer := myHttp.New(handler, conf.Port)

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
