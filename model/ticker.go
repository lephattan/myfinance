package model

import (
	"database/sql"
	"fmt"
	"myfinance/database"
	"myfinance/helper"
	"net/url"
	"strings"
)

const TickerTablename = "tickers"

type Ticker struct {
	Symbol         string                   `db:"symbol" form:"symbol"`
	Name           string                   `db:"name" form:"name"`
	CurrentPrice   database.Nullable[int64] `db:"current_price" form:"current_price"`
	PriceUpdatedAt database.Nullable[int64] `db:"price_updated_at"`
}

func (t *Ticker) TableName() string {
	return TickerTablename
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
	return helper.ModelScan(t, rows)
}

func (t *Ticker) String() string {
	return fmt.Sprintf("<Ticker Symbol: %s, Name: %s>", t.Symbol, t.Name)
}

// List of Tickers
type Tickers []*Ticker

func (t *Tickers) TableName() string {
	return "tickers"
}

func (t *Tickers) SearchColumns() []string {
	return []string{"symbol", "name"}
}

// Scan binds mysql rows to this Tickers. NOTE: wtf is this
func (ts *Tickers) Scan(rows *sql.Rows) (err error) {
	err = helper.ModelListScan(ts, rows)
	return
}

func (t *Tickers) ParseListOptions(q *url.Values) database.ListOptions {
	opt := database.ParseListOptions(q)
	opt.SetTableName(t.TableName())
	if search := q.Get("s"); search != "" {
		opt.Search(t.SearchColumns(), search)
	}

	if symbol := q.Get("symbol"); symbol != "" {
		opt.Where("symbol", symbol)
	}

	return opt
}

// The requests.
// Ref: https://github.com/kataras/iris/blob/24ba4e8933b9d58096a56e5c6f2de968f80eb602/_examples/database/mysql/entity/product.go#L73C10-L73C10
type (
	CreateTickerRequest struct { // all required.
		Symbol string `json:"symbol"`
	}

	TickerUpdate struct {
		Symbol       string                   `db:"symbol"`
		Name         string                   `db:"name" form:"name"`
		CurrentPrice database.Nullable[int64] `db:"current_price" form:"current-price"`
	}

	GetTickerRequest struct {
		Symbol string `json:"symbol"` // required
	}

	DeleteTickerRequest GetTickerRequest

	GetTickersRequest struct {
		Symbol     string `query:"symbol"`
		Name       string `query:"name"`
		Search     string `query:"s"`
		Pagination *Pagination
	}
)

func NewGetTickersRequest() GetTickersRequest {
	pagination := NewPagination()
	query := GetTickersRequest{
		Pagination: &pagination,
	}
	return query
}

func (t *TickerUpdate) UpdateQuery() (query string, args []interface{}) {
	query = fmt.Sprintf(`Update %s
		Set 
			name = ?,
			current_price = ?
		Where %s = Lower(?);`,
		TickerTablename,
		"symbol",
	)
	args = append(args, &t.Name, &t.CurrentPrice, &t.Symbol)
	return
}

func (t *TickerUpdate) ValidateUpdate() bool {
	return true
}
