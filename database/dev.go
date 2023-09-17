package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

type devdb struct{}

/* Execute query string */
func (db *devdb) Exec(q string) error {
	fmt.Println("query: " + q)
	return nil
}

/*
Driver for database migration
Return:
- (string) database type
- (database.Driver) database driver
- (error) error when creating database driver
*/
func (db *devdb) MigrateDriver() (string, database.Driver, error) {
	conn, err := sql.Open("sqlite3", filepath.Join(".", "dev-database.db"))
	if err != nil {
		panic("Error opening dev db connection: " + err.Error())
	}
	driver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	return "sqlite3", driver, err
}
