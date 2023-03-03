package database

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
)

type Config struct {
	DatabaseUser     string
	DatabaseName     string
	DatabaseHost     string
	DatabasePassword string
}

type PostgreSQL struct {
	DB *db.Session
}

// NewDbSess create session with PostgreSQL.
func NewDbSess(conf *Config) db.Session {
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
