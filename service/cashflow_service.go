package service

import (
	"context"
	"database/sql"
	"fmt"
	"myfinance/database"
	"myfinance/helper"
	"myfinance/model"
)

type CashflowService interface {
	List(ctx context.Context, dest interface{}) error
}

func NewCasflowService(db database.DB) (service CashflowService) {
	service = &cashflow{db: db}
	return
}

type cashflow struct {
	db database.DB
}

func (c *cashflow) List(ctx context.Context, dest interface{}) error {
	query := fmt.Sprintf(`SELECT t.date,
    SUM(CASE 
        When t.transaction_type = 'sell' THEN 
            t.volume * t.price - commission
        END) inflow,
    SUM(CASE 
        When t.transaction_type = 'buy' THEN 
            t.volume * t.price + commission 
        END) outflow
		From 
				%s as t
		Group By t.date;`, model.TransactionsTablename)
	args := []interface{}{}
	rows, err := c.db.Select(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return helper.ModelListScan(dest, rows)
}
