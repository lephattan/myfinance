package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ListOptions struct {
	Table         string // the table name.
	Offset        uint64 // inclusive.
	Limit         uint64
	OrderByColumn string
	Order         string // "ASC" or "DESC" (could be a bool type instead).
	WhereColumn   string
	WhereValue    interface{}
	SearchFor     string
	SearchColumns []string
}

func (opt *ListOptions) Where(col_name string, col_value interface{}) {
	opt.WhereColumn = col_name
	opt.WhereValue = col_value
}

func (opt *ListOptions) Search(searchCols []string, searchFor string) {
	opt.SearchColumns = searchCols
	opt.SearchFor = searchFor
}

func (opt *ListOptions) SetTableName(name string) {
	opt.Table = name
}

func (opt ListOptions) BuildQuery() (q string, args []interface{}) {
	q = fmt.Sprintf("SELECT * From %s", opt.Table)

	if opt.WhereColumn != "" && opt.WhereValue != nil {
		q += fmt.Sprintf(" Where %s = :whereValue", opt.WhereColumn)
		args = append(args, sql.Named("whereValue", opt.WhereValue))
	} else if len(opt.SearchColumns) != 0 && opt.SearchFor != "" {
		whereClaudes := []string{}
		for _, col := range opt.SearchColumns {
			whereClaudes = append(whereClaudes, fmt.Sprintf("%s Like :search", col))
		}
		q += fmt.Sprintf(" Where %s", strings.Join(whereClaudes, " Or "))
		args = append(args, sql.Named("search", "%"+opt.SearchFor+"%"))
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
	q += ";"
	return
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

const DefaultLimit = 20

// ParseListOptions returns a `ListOptions` from a map[string][]string.
func ParseListOptions(q *url.Values) ListOptions {
	page, _ := strconv.ParseUint(q.Get("page"), 10, 64)
	limit, _ := strconv.ParseUint(q.Get("per_page"), 10, 64)
	if limit == 0 {
		limit = DefaultLimit
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	order := q.Get("order")

	return ListOptions{Offset: offset, Limit: limit, Order: order}
}
