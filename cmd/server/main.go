package main

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/config/container"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/database"
	http2 "github.com/BogdanStaziyev/softcery-test/internal/infra/transport"
	"github.com/pkg/errors"
	"log"
	"os"
)

func main() {
	//initialize configuration
	var conf = config.GetConfiguration()

	//make migration
	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	//initialize storage location
	_, err = os.Stat(conf.FileStorageLocation)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(conf.FileStorageLocation, os.ModePerm)
		if err != nil {
			log.Fatalf("storage folder can not be created %s", err)
		}
	} else if err != nil {
		log.Fatalf("storage folder is not available %s", err)
	}

	//initialize container with handlers services and db
	cont := container.New(conf)

	//initialize server
	srv := http2.NewServer()

	//start router
	http2.EchoRouter(srv, cont)

	//start server
	err = srv.Start()
	if err != nil {
		log.Fatal("Port already used")
	}
}
