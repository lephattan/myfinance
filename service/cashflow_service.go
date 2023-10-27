package service

import (
	"context"
	"database/sql"
	"myfinance/database"
	"myfinance/helper"
	"myfinance/model"
)

type CashflowService interface {
	List(ctx context.Context, opt model.CashflowListingOptions, dest interface{}) error
}

func NewCasflowService(db database.DB) (service CashflowService) {
	service = &cashflow{db: db}
	return
}

type cashflow struct {
	db database.DB
}

func (c *cashflow) List(ctx context.Context, opt model.CashflowListingOptions, dest interface{}) error {
	query, args := opt.PrepareQuery()
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
