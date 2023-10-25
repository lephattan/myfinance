package model

import (
	"database/sql"
	"fmt"
	"myfinance/database"
	"myfinance/helper"
	"net/url"
	"strings"
)

const PortfolioTablename = "portfolios"

type Portfolio struct {
	ID          uint64 `db:"id,omitinsert" json:"id"`
	Name        string `db:"name" json:"name" form:"name"`
	Description string `db:"description" json:"description" form:"description"`
}

func (t *Portfolio) TableName() string {
	return PortfolioTablename
}

func (t *Portfolio) PrimaryKey() string {
	return "id"
}

func (t *Portfolio) SortBy() string {
	return "id"
}

func (t *Portfolio) ValidateInsert() bool {
	return strings.TrimSpace(t.Name) != ""
}

func (t *Portfolio) Scan(rows *sql.Rows) error {
	return rows.Scan(&t.ID, &t.Name, &t.Description)
}

func (t *Portfolio) String() string {
	return fmt.Sprintf("<Portfolio ID: %d, Name: %s>", t.ID, t.Name)
}

// List of Portfolios
type Portfolios []*Portfolio

// Scan binds mysql rows to this Portfolios
func (ts *Portfolios) Scan(rows *sql.Rows) (err error) {
	return helper.ModelListScan(ts, rows)
}

func (t *Portfolios) SearchColumns() []string {
	return []string{"id", "name"}
}

func (t *Portfolios) ParseListOptions(q *url.Values) database.ListOptions {
	opt := database.ParseListOptions(q)
	opt.SetTableName(PortfolioTablename)
	if search := q.Get("s"); search != "" {
		opt.Search(t.SearchColumns(), search)
	}

	if name := q.Get("name"); name != "" {
		opt.Where("name", name)
	}

	if id := q.Get("id"); id != "" {
		opt.Where("id", id)
	}

	return opt
}

// the requests

type (
	GetPortfoliosRequest struct {
		Name       string `query:"name"`
		Search     string `query:"s"`
		Pagination *Pagination
	}
)

func NewGetPortfoliosRequest() GetPortfoliosRequest {
	pagination := NewPagination()
	query := GetPortfoliosRequest{
		Pagination: &pagination,
	}
	return query
}

type HoldingSummarry struct {
	TotalCost  int64
	TotalValue int64
}

func (h *HoldingSummarry) GainLoss() int64 {
	return h.TotalValue - h.TotalCost
}

func (h *HoldingSummarry) GainLossPercent() float64 {
	if h.TotalCost == 0 {
		return 0
	}
	gain_loss := h.GainLoss()
	percent := float64(gain_loss) / float64(h.TotalCost) * 100
	return percent
}
