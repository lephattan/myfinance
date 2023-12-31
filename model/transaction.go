package model

import (
	"database/sql"
	"errors"
	"fmt"
	"myfinance/database"
	"myfinance/helper"
	"net/url"
	"strings"
)

type TransactionType string

const (
	Buy  TransactionType = "buy"
	Sell                 = "sell"
)

const TransactionsTablename = "transactions"

// transaction table name
func TransactionsTable() string { return TransactionsTablename }

// transaction columns to perform search query
func TransactionsSearchColumns() []string {
	return []string{
		"date",
		"ticker_symbol",
		"transaction_type",
		"note",
	}
}

var TransactionTypes = [...]TransactionType{
	Buy,
	Sell,
}

type Transaction struct {
	ID              uint64             `db:"id,omitinsert"`
	Date            database.UnixDate  `db:"date" form:"date"`
	TickerSymbol    string             `db:"ticker_symbol" form:"ticker-symbol"`
	PortfolioID     uint64             `db:"portfolio_id" form:"portfolio-id"`
	TransactionType TransactionType    `db:"transaction_type" form:"transaction-type"`
	Volume          uint64             `db:"volume" form:"volume"`
	Price           uint64             `db:"price" form:"price"`
	Commission      uint64             `db:"commission" form:"commission"`
	Note            sql.NullString     `db:"note" form:"note"`
	RefID           database.NullInt64 `db:"ref_id" form:"ref-id"`
}

func (t Transaction) TableName() string {
	return TransactionsTablename
}

func (t Transaction) PrimaryKey() string {
	return "id"
}

func (t Transaction) SortBy() string {
	return "id"
}

func (t *Transaction) ValidateInsert() error {
	if t.Date.Init64() <= 0 {
		return errors.New("invalid transaction date")
	}
	if t.TickerSymbol == "" {
		return errors.New("missing transaction ticker symbol")
	}
	if err := t.ValidateType(); err != nil {
		return err
	}
	if t.Volume < 0 {
		return errors.New("negative transaction volume")
	}
	if t.PortfolioID < 0 {
		return errors.New("invalid transaction portfolio_id")
	}
	return nil
}

// Validate TransactionType
func (t *Transaction) ValidateType() error {
	for _, trans_type := range TransactionTypes {
		if t.TransactionType == trans_type {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("transaction type: %v is not implemented", t.TransactionType))
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

func (t *Transaction) Total() (total int64, err error) {
	switch t.TransactionType {
	case "buy":
		{
			return int64(t.Price*t.Volume + t.Commission), nil
		}
	case "sell":
		{
			return int64(t.Price*t.Volume - t.Commission), nil
		}
	default:
		return 0, errors.New(fmt.Sprintf("transaction total is not implemented for type: %s", t.TransactionType))
	}

}

// List of Transactions
type Transactions []*Transaction

// Scan binds mysql rows to list of Transactions
func (ts *Transactions) Scan(rows *sql.Rows) (err error) {
	return helper.ModelListScan(ts, rows)
}

func (t *Transactions) ParseListOptions(q *url.Values) database.ListOptions {
	opt := database.ParseListOptions(q)
	opt.SetTableName(TransactionsTable())
	if search := q.Get("s"); search != "" {
		opt.Search(TransactionsSearchColumns(), search)
	}

	if id := q.Get("id"); id != "" {
		opt.Where("id", id)
	}

	if date := q.Get("date"); date != "" {
		opt.Where("date", date)
	}

	if portfolio_id := q.Get("portfolio_id"); portfolio_id != "" {
		opt.Where("portfolio_id", portfolio_id)
	}

	if transaction_type := q.Get("transaction_type"); transaction_type != "" {
		opt.Where("transaction_type", transaction_type)
	}

	if opt.OrderByColumn == "" {
		opt.OrderByColumn = "date"
	}
	if opt.Order == "" {
		opt.Order = "DESC"
	}

	return opt
}
