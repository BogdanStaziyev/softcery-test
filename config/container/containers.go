package container

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
)

type Container struct {
	Services
	Controllers
}

type Services struct {
}

type Controllers struct {
}

func New(conf config.Configuration) Container {
	_ = getDbSess(conf)

	return Container{
		Services:    Services{},
		Controllers: Controllers{},
	}
}

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
