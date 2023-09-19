package database

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

type devdb struct{}

const DB_TYPE = "sqlite3"

/*
Driver for database migration
Return:
- (string) database type
- (database.Driver) database driver
- (error) error when creating database driver
*/
func (db *devdb) MigrateDriver() (string, database.Driver, error) {
	conn, err := sql.Open(DB_TYPE, filepath.Join(".", "dev-database.db"))
	if err != nil {
		panic("Error opening dev db connection: " + err.Error())
	}
	driver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	return DB_TYPE, driver, err
}

type dbconn struct {
	dsn string
}

func (db *devdb) Connect(s dbconn) *sql.DB {
	if s.dsn == "" {
		s.dsn = "file:dev-database.db?cache=shared"
	}
	conn, err := sql.Open(DB_TYPE, s.dsn)
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}
	return conn
}

/* Execute query string and return result*/
func (db *devdb) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	conn := db.Connect(dbconn{})
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		panic("Error querying database: " + err.Error())
	}
	defer rows.Close()
	if scannable, ok := dest.(Scannable); ok {
		return scannable.Scan(rows)
	}
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return rows.Scan(dest)
}

/* Similar to Select but does not return the result*/
func (db *devdb) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	conn := db.Connect(dbconn{})
	defer conn.Close()
	res, err := conn.ExecContext(ctx, query, args...)
	return res, err
}

func (db *devdb) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	conn := db.Connect(dbconn{})
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}

	if Scannable, ok := dest.(Scannable); ok {
		return Scannable.Scan(rows)
	}
	return rows.Scan(dest)
}
