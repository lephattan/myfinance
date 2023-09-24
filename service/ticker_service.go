package service

import (
	"context"
	"fmt"
	"myfinace/database"
	"myfinace/env"
	"myfinace/model"
	"strings"
)

type TickerService interface {
	List(ctx context.Context, opt database.ListOptions, dest interface{}) error
	Create(ctx context.Context, t model.Ticker) (int64, error)
	Get(ctx context.Context, symbol string, dest interface{}) (err error)
	Update(ctx context.Context, t model.Ticker) (int, error)
}

func NewTickerService(e env.Env, db database.DB) TickerService {
	service := &ticker{db: db}
	return service
}

type ticker struct {
	db  database.DB
	rec database.Record
}

func (s *ticker) List(ctx context.Context, opt database.ListOptions, dest interface{}) error {
	q, args := opt.BuildQuery()
	err := s.db.Select(ctx, dest, q, args...)
	return err
}

func (s *ticker) Create(ctx context.Context, t model.Ticker) (int64, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	query := fmt.Sprintf("Insert Into %s (symbol, name) Values (Lower(?),?);", t.TableName())
	res, err := s.db.Exec(ctx, query, strings.TrimSpace(t.Symbol), t.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *ticker) Get(ctx context.Context, symbol string, dest interface{}) (err error) {
	q := fmt.Sprintf("Select * From %s Where `symbol` = ?;", "tickers")
	err = s.db.Get(ctx, dest, q, strings.TrimSpace(symbol))
	return
}

func (s *ticker) Update(ctx context.Context, t model.Ticker) (int, error) {
	if !t.ValidateInsert() {
		return 0, database.ErrUnprocessable
	}
	q := fmt.Sprintf(`Update %s
		Set
			name = ?
		Where %s = Lower(?);
		`, t.TableName(), t.PrimaryKey())
	res, err := s.db.Exec(ctx, q, t.Name, strings.ToLower(t.Symbol))
	if err != nil {
		return 0, err
	}
	n := database.GetAffectedRows(res)
	return n, err
}
