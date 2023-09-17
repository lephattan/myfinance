package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4/database"
)

type proddb struct{}

func (db *proddb) Exec(q string) error {
	return fmt.Errorf("production database is not implemented")
}

func (db *proddb) MigrateDriver() (string, database.Driver, error) {
	panic("not implemented")
}
