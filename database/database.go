package database

import (
	"myfinace/env"

	"github.com/golang-migrate/migrate/v4/database"
)

type DB interface {
	Exec(q string) error
	MigrateDriver() (string, database.Driver, error)
}

func NewDB(e env.Env) DB {
	switch e {
	case env.PROD:
		return &proddb{}
	case env.DEV:
		return &devdb{}
	case env.TESTING:
		return &devdb{}
	default:
		panic("unknown environment")
	}
}
