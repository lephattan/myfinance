package database

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4/database"
)

type proddb struct{}

func (db *proddb) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	panic("not implemented")
}

func (db *proddb) MigrateDriver() (string, database.Driver, error) {
	panic("not implemented")
}
func (db *proddb) Select(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	panic("not implemented")
}
func (db *proddb) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	panic("not implemented")
}
