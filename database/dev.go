package database

import "fmt"

type devdb struct{}

func (db *devdb) Exec(q string) error {
	fmt.Println("query: " + q)
	return nil
}
