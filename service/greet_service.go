package service

import (
	"fmt"
	"myfinace/database"
	"myfinace/env"
)

type GreetSerivce interface {
	Say(input string) (string, error)
}

func NewGreetService(e env.Env, db database.DB) GreetSerivce {
	service := &greeter{db: db, prefix: "Hello"}
	switch e {
	case env.PROD:
		return service
	case env.DEV, env.TESTING:
		return &greeterWithLogging{service}
	default:
		panic("unknown env")
	}
}

type greeter struct {
	prefix string
	db     database.DB
}

func (s *greeter) Say(input string) (string, error) {
	if err := s.db.Exec("simulate query..."); err != nil {
		return "", err
	}
	result := s.prefix + " " + input
	return result, nil
}

type greeterWithLogging struct {
	*greeter
}

func (s *greeterWithLogging) Say(input string) (string, error) {
	result, err := s.greeter.Say(input)
	fmt.Printf("result: %s\nerror:%v\n", result, err)
	return result, err
}
