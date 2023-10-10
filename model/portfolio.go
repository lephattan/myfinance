package model

import (
	"database/sql"
	"fmt"
	"myfinance/helper"
	"strings"
)

type Portfolio struct {
	ID          uint64 `db:"id,omitinsert" json:"id"`
	Name        string `db:"name" json:"name" form:"name"`
	Description string `db:"description" json:"description" form:"description"`
}

func (t *Portfolio) TableName() string {
	return "portfolios"
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
