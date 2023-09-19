package model

import (
	"database/sql"
	"fmt"
	"strings"
)

type Ticker struct {
	Symbol string `db:"symbol" json:"symbol"`
	Name   string `db:"name" json:"name"`
}

func (t *Ticker) TableName() string {
	return "tickers"
}

func (t *Ticker) PrimaryKey() string {
	return "symbol"
}

func (t *Ticker) SortBy() string {
	return "symbol"
}

func (t *Ticker) ValidateInsert() bool {
	return strings.TrimSpace(t.Symbol) != ""
}

func (t *Ticker) Scan(rows *sql.Rows) error {
	return rows.Scan(&t.Symbol, &t.Name)
}

func (t *Ticker) String() string {
	return fmt.Sprintf("<Ticker Symbol: %s, Name: %s>", t.Symbol, t.Name)
}

// List of Tickers
type Tickers []*Ticker

// Scan binds mysql rows to this Categories. NOTE: wtf is this
func (ts *Tickers) Scan(rows *sql.Rows) (err error) {
	cp := *ts
	for rows.Next() {
		t := new(Ticker)
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

// The requests.
// Ref: https://github.com/kataras/iris/blob/24ba4e8933b9d58096a56e5c6f2de968f80eb602/_examples/database/mysql/entity/product.go#L73C10-L73C10
type (
	CreateTickerRequest struct { // all required.
		Symbol string `json:"symbol"`
	}

	UpdateProductRequest struct {
		Name string `json:"name"`
	} // at least 1 required.

	GetTickerRequest struct {
		Symbol string `json:"symbol"` // required
	}

	DeleteTickerRequest GetTickerRequest

	// GetTickersRequest struct {
	// 	// [page, offset...]
	// }
)
