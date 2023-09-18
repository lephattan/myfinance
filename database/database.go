package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myfinace/env"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
)

type DB interface {
	Exec(q string) error
	MigrateDriver() (string, database.Driver, error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Scannable interface {
	Scan(*sql.Rows) error
}

// Record should represent a database record.
// It holds the table name and the primary key.
// Entities should implement that
// in order to use the BaseService's methods.
type Record interface {
	TableName() string  // the table name which record belongs to.
	PrimaryKey() string // the primary key of the record.
}

func NewDB(e env.Env) DB {
	switch e {
	case env.PROD:
		return &proddb{}
	case env.DEV:
		return &devdb{}
	case env.TESTING:
		return &devdb{}
	default:
		panic("unknown environment")
	}
}

type ListOptions struct {
	Table         string // the table name.
	Offset        uint64 // inclusive.
	Limit         uint64
	OrderByColumn string
	Order         string // "ASC" or "DESC" (could be a bool type instead).
	WhereColumn   string
	WhereValue    interface{}
}

func (opt ListOptions) Where(col_name string, col_value interface{}) ListOptions {
	opt.WhereColumn = col_name
	opt.WhereValue = col_value
	return opt
}

func (opt ListOptions) BuildQuery() (q string, args []interface{}) {
	q = fmt.Sprintf("SELECT * From %s", opt.Table)

	if opt.WhereColumn != "" && opt.WhereValue != nil {
		q += fmt.Sprintf(" Where %s = ?", opt.WhereColumn)
		args = append(args, opt.WhereValue)
	}

	if opt.OrderByColumn != "" {
		q += fmt.Sprintf(" Order By %s %s", opt.OrderByColumn, ParseOrder(opt.Order))
	}

	if opt.Limit > 0 {
		q += fmt.Sprintf(" Limit %d", opt.Limit)
	}

	if opt.Offset > 0 {
		q += fmt.Sprintf(" Offset %d", opt.Offset)
	}
	return
}

func (opt ListOptions) TableName(name string) {
	opt.Table = name
}

const (
	ascending  = "ASC"
	descending = "DESC"
)

func ParseOrder(order string) string {
	order = strings.TrimSpace(order)
	if len(order) > 4 {
		if strings.HasPrefix(strings.ToUpper(order), descending) {
			return descending
		}
	}
	return ascending
}

const defaultLimit = 20

// ParseListOptions returns a `ListOptions` from a map[string][]string.
func ParseListOptions(q url.Values) ListOptions {
	page, _ := strconv.ParseUint(q.Get("page"), 10, 64)
	limit, _ := strconv.ParseUint(q.Get("per_page"), 10, 64)
	if limit == 0 {
		limit = defaultLimit
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	order := q.Get("order")

	return ListOptions{Offset: offset, Limit: limit, Order: order}
}

var ErrUnprocessable = errors.New("invalid model")
