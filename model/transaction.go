package model

import (
	"database/sql"
	"fmt"
	"log"
	"myfinace/database"
	"strings"
)

type TransactionType string

const (
	Buy  TransactionType = "buy"
	Sell                 = "sell"
)

var TransactionTypes = [...]TransactionType{
	Buy,
	Sell,
}

type Transaction struct {
	ID              uint64             `db:"id"`
	Date            uint64             `db:"date"`
	TickerSymbol    string             `db:"ticker_symbol"`
	PortfolioID     uint64             `db:"portfolio_id"`
	TransactionType TransactionType    `db:"transaction_type"`
	Volume          uint64             `db:"volume"`
	Price           uint64             `db:"price"`
	Commission      uint64             `db:"commission"`
	Note            sql.NullString     `db:"note"`
	RefID           database.NullInt64 `db:"ref_id"`
}

func (t Transaction) TableName() string {
	return "transactions"
}

func (t Transaction) PrimaryKey() string {
	return "id"
}

func (t Transaction) SortBy() string {
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
	for _, trans_type := range TransactionTypes {
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
		&t.PortfolioID,
		&t.TransactionType,
		&t.Volume,
		&t.Price,
		&t.Commission,
		&t.Note,
		&t.RefID,
	)
}

// Preresentative string of a transaction
func (t *Transaction) String() string {
	return fmt.Sprintf(
		"<Transaction %s %d %s @%d in PortfolioID: %d>",
		t.TransactionType,
		t.Volume,
		strings.ToUpper(t.TickerSymbol),
		t.Price,
		t.PortfolioID,
	)
}

func (t *Transaction) GenerateInsertStatement() (stmt string, args []interface{}, err error) {
	return database.GenerateInsertStatement(*t)
}

// List of Transactions
type Transactions []*Transaction

// Scan binds mysql rows to list of Transactions
func (ts *Transactions) Scan(rows *sql.Rows) (err error) {
	cp := *ts
	for rows.Next() {
		t := new(Transaction)
		if err = t.Scan(rows); err != nil {
			log.Printf("Error scanning row: %s", err.Error())
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
