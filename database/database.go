package database

import "myfinace/env"

type DB interface {
	Exec(q string) error
}

func NewDB(e env.Env) DB {
	switch e {
	case env.PROD:
		return &proddb{}
	case env.DEV:
		return &devdb{}
	default:
		panic("unknown environment")
	}
}
