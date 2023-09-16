package database

import "fmt"

type proddb struct{}

func (db *proddb) Exec(q string) error {
	return fmt.Errorf("production database is not implemented")
}
