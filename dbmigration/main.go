package main

import (
	"myfinance/database"
	"myfinance/env"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	env := env.ReadEnv("APP_ENV", env.PROD)
	db := database.NewDB(env)
	db_type, driver, err := db.MigrateDriver()
	if err != nil {
		panic("Error creating migration driver: " + err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		db_type,
		driver,
	)
	m.Up()
}
