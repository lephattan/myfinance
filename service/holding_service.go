package service

import "myfinance/database"

type HoldingService interface {
}

func NewHoldingService(db database.DB) HoldingService {
	service := &holding{db: db}
	return service
}

type holding struct {
	db database.DB
}
