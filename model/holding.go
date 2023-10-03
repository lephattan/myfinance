package model

import (
	"database/sql"
	"myfinance/database"
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

func (h Holding) TableName() string {
	return "holdings"
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
	return ModelScan(h, rows)
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
