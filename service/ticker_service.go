package service

import (
	"context"
	"fmt"
	"myfinace/database"
	"myfinace/env"
	"myfinace/model"
)

type TickerService interface {
	List(ctx context.Context, dest interface{}) error
	Create(ctx context.Context, t model.Ticker) (int64, error)
}

func NewTickerService(e env.Env, db database.DB) TickerService {
	service := &ticker{db: db}
	return service
}

type ticker struct {
	db  database.DB
	rec database.Record
}

func (s *ticker) List(ctx context.Context, dest interface{}) error {
	opt := database.ListOptions{
		Table: "tickers",
	}
	q, args := opt.BuildQuery()
	err := s.db.Select(ctx, dest, q, args...)
	return err
}

func (s *ticker) Create(ctx context.Context, t model.Ticker) (int64, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	query := fmt.Sprintf("Insert Into %s (symbol, name) Values (Lower(?),?);", t.TableName())
	res, err := s.db.Exec(ctx, query, t.Symbol, t.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
