package model

import (
	"math"
	"myfinance/database"
)

type Pagination struct {
	Page    uint64 `query:"page"`
	PerPage uint64 `query:"per_page"`
	Count   uint64
}

func (p *Pagination) MaxPage() uint64 {
	max_page := math.Ceil(float64(p.Count) / float64(p.PerPage))
	return uint64(max_page)
}

// Set default value for Pagination
func (p *Pagination) Default() {
	p.Page = 1
	p.PerPage = database.DefaultLimit
	p.Count = 0
}

// Create new Pagination object with default value
func NewPagination() Pagination {
	var pagination Pagination
	pagination.Default()
	return pagination
}
