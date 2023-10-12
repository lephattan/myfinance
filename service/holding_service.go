package service

import (
	"context"
	// "database/sql"
	// "fmt"
	"myfinance/database"
	"myfinance/helper"
	"myfinance/model"
)

type HoldingService interface {
	Create(ctx context.Context, h model.Holding) error
	List(ctx context.Context, opt database.ListOptions, dest interface{}) error
	Get(ctx context.Context, opt database.ListOptions, dest interface{}) error
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

func (s *holding) List(ctx context.Context, opt database.ListOptions, dest interface{}) (err error) {
	opt.SetTableName(model.HoldingTablename)
	q, args := opt.BuildQuery()
	rows, err := s.db.Select(ctx, q, args...)
	if err != nil {
		return err
	}
	return helper.ModelListScan(dest, rows)
}

func (s *holding) Get(ctx context.Context, opt database.ListOptions, dest interface{}) (err error) {
	opt.SetTableName(model.HoldingTablename)
	opt.Limit = 1
	q, args := opt.BuildQuery()
	rows, err := s.db.Select(ctx, q, args...)
	if err != nil {
		return err
	}
	return helper.ModelScan(dest, rows)
}
