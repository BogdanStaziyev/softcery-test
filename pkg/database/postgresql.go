// Package database implements postgres connection.
package database

import (
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type PostgreSQL struct {
	session db.Session
	Coll    db.Collection
	Builder squirrel.StatementBuilderType
}

type pSQL interface {
	Collection(s string) db.Collection
}

var _ pSQL = (*PostgreSQL)(nil)

// NewDbSess create session with PostgreSQL.
func NewDbSess(conf *Config) *PostgreSQL {
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
	return &PostgreSQL{
		session: sess,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Colon)}
}

func (p *PostgreSQL) Collection(s string) db.Collection {
	return p.session.Collection(s)
}

type Config struct {
	DatabaseUser     string
	DatabaseName     string
	DatabaseHost     string
	DatabasePassword string
}
