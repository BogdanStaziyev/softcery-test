package main

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/app"
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
)

func main() {
	l := logger.New("main")
	//initialize configuration
	var conf = config.GetConfiguration()
	l.Info("Success read config")

	//run application
	app.Run(conf)
}
