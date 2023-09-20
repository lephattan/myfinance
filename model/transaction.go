package model

import (
	"database/sql"
	"fmt"
	"strings"
)

type TransactionType string

const (
	Buy  TransactionType = "buy"
	Sell                 = "sell"
)

type Transaction struct {
	ID              uint64          `db:"id"`
	Date            uint64          `db:"date"`
	TickerSymbol    string          `db:"ticker_symbol"`
	TransactionType TransactionType `db:"transaction_type"`
	Volume          uint64          `db:"volume"`
	Price           uint64          `db:"price"`
	Commission      uint64          `db:"commission"`
	Note            string          `db:"note"`
	PortfolioID     uint64          `db:"portfolio_id"`
	RefID           uint64          `db:"ref_id"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}

func (t *Transaction) PrimaryKey() string {
	return "id"
}

func (t *Transaction) SortBy() string {
	return "id"
}

func (t *Transaction) ValidateInsert() bool {
	return t.Date > 0 &&
		t.TickerSymbol != "" &&
		t.ValidateType() &&
		t.Volume > 0 &&
		t.PortfolioID > 0
}

// Validate TransactionType
func (t *Transaction) ValidateType() bool {
	types := []TransactionType{
		Buy,
		Sell,
	}
	for _, trans_type := range types {
		if t.TransactionType == trans_type {
			return true
		}
	}
	return false
}

// Binds sql rows to transaction
func (t *Transaction) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&t.ID,
		&t.Date,
		&t.TickerSymbol,
		&t.TransactionType,
		&t.Volume,
		&t.Price,
		&t.Commission,
		&t.Note,
		&t.PortfolioID,
		&t.RefID,
	)
}

// Preresentative string of a transaction
func (t *Transaction) String() string {
	return fmt.Sprintf(
		"<Transaction %s %s @%d of PortfolioID: %d>",
		t.TransactionType,
		strings.ToUpper(t.TickerSymbol),
		t.Price,
		t.PortfolioID,
	)
}

// List of Transactions
type Transactions []*Transaction

// Scan binds mysql rows to list of Transactions
func (ts *Transactions) Scan(rows *sql.Rows) (err error) {
	cp := *ts
	for rows.Next() {
		t := new(Transaction)
		if err = t.Scan(rows); err != nil {
			return
		}
		cp = append(cp, t)
	}

	if len(cp) == 0 {
		return sql.ErrNoRows
	}
	*ts = cp
	return rows.Err()

}
