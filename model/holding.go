package model

import (
	"database/sql"
	"errors"
	"fmt"
	"myfinance/database"
	"myfinance/helper"
)

type Holding struct {
	Symbol       string                   `db:"symbol"`
	PortfolioID  uint64                   `db:"portfolio_id"`
	TotalShares  int64                    `db:"total_shares"`
	TotalCost    int64                    `db:"total_cost"`
	AveragePrice int64                    `db:"average_price"`
	CurrentValue database.Nullable[int64] `db:"current_value"`
	UpdatedAt    database.Nullable[int64] `db:"updated_at"`
}

const HoldingTablename = "holdings"

func (h Holding) TableName() string {
	return HoldingTablename
}

func (h Holding) PrimaryKey() string {
	return ""
}

func (h Holding) SortBy() string {
	return "symbol"
}

func (h *Holding) ValidateInsert() bool {
	return true
}

func (h *Holding) Scan(rows *sql.Rows) error {
	return helper.ModelScan(h, rows)
}

func (h *Holding) GenerateInsertStatement() (stmt string, args []interface{}, err error) {
	return database.GenerateInsertStatement(*h)
}

func (h *Holding) GetAveragePrice() int64 {
	if h.TotalShares == 0 {
		return 0
	}
	return int64(h.TotalCost / h.TotalShares)
}

func (h *Holding) HandleTransaction(t *Transaction) (err error) {
	switch t.TransactionType {
	case "buy":
		{
			h.TotalShares += int64(t.Volume)
			total, err := t.Total()
			if err != nil {
				return err
			}
			h.TotalCost += total
		}
	case "sell":
		{
			h.TotalShares -= int64(t.Volume)
			total, err := t.Total()
			if err != nil {
				return err
			}
			h.TotalCost -= total
		}
	default:
		{
			return errors.New(
				fmt.Sprintf("cannot calculate holding from transaction type %s", t.TransactionType),
			)
		}
	}
	return nil
}

func (h *Holding) GainLoss() int64 {
	if h.CurrentValue.Valid {
		return h.CurrentValue.Actual - h.TotalCost
	}
	return 0
}

func (h *Holding) GainLossPercent() float64 {
	if h.TotalCost == 0 {
		return 0
	}
	gl := float64(h.GainLoss()) / float64(h.TotalCost)
	return gl * 100
}
