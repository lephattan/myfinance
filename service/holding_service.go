package service

import (
	"context"
	// "database/sql"
	// "fmt"
	"myfinance/database"
	"myfinance/model"
)

type HoldingService interface {
	Create(ctx context.Context, h model.Holding) error
}

func NewHoldingService(db database.DB) HoldingService {
	service := &holding{db: db}
	return service
}

type holding struct {
	db database.DB
}

func (s *holding) Create(ctx context.Context, h model.Holding) error {
	if !h.ValidateInsert() {
		return database.ErrUnprocessable
	}
	query, args, err := h.GenerateInsertStatement()
	_, err = s.db.Exec(
		ctx,
		query,
		args...,
	)
	return err
}
