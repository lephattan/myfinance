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
	fmt.Printf("%+v\n", opt)
	q, args := opt.BuildQuery()
	err := s.db.Select(ctx, dest, q, args...)
	return err
}

func (s *ticker) Create(ctx context.Context, t model.Ticker) (model.Ticker, error) {
	if !t.ValidateInsert() {
		return model.Ticker{}, database.ErrUnprocessable
	}
	query := fmt.Sprintf("Insert Into %s (symbol, name) Value (?,?);", t.TableName())
	newTicker := new(model.Ticker)
	err := s.db.Select(ctx, newTicker, query, t.Symbol, t.Name)
	return *newTicker, err
}
