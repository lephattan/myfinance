package main

import (
	"log"
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
	version, dirty, err := m.Version()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Current version: %d, dirty: %v", version, dirty)

	err = m.Up()
	if err == nil {
		log.Println("Migration is success")
		return
	}
	switch err {
	case migrate.ErrNoChange:
		log.Println(err)
	default:
		log.Panic(err)
	}

}
