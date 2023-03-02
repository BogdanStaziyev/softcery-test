package main

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/app"
	"log"
)

func main() {
	//initialize configuration
	var conf = config.GetConfiguration()
	log.Println("Success read config")

	app.Run(conf)
}
