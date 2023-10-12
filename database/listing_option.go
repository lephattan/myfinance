package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type WhereGroup struct {
	WhereClauses []*WhereClause
	// AND or OR operator for where clause
	Operator string
}

func (w *WhereGroup) Where(where_column string, where_value interface{}, operator string) {
	where_clause := WhereClause{
		WhereColumn: where_column,
		WhereValue:  where_value,
		Operator:    operator,
	}
	w.WhereClauses = append(w.WhereClauses, &where_clause)
}

func (w *WhereGroup) BuildQuery() (query string, args []interface{}) {
	for i, where_clause := range w.WhereClauses {
		if i != 0 {
			query += " " + strings.ToUpper(where_clause.Operator)
		}
		where_query, where_args := where_clause.BuildQuery()
		query += " " + where_query
		args = append(args, where_args...)
	}
	if len(w.WhereClauses) > 1 {
		query = "(" + query + ")"
	}
	return
}

// Where clause for sql query, the Operator of first Where clause will be ignored
type WhereClause struct {
	WhereColumn string
	WhereValue  interface{}
	// AND or OR operator for where clause
	Operator string
}

func (w *WhereClause) BuildQuery() (query string, args []interface{}) {
	query = fmt.Sprintf("%s = ?", w.WhereColumn)
	args = append(args, w.WhereValue)
	return
}

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
	WhereGroups   []*WhereGroup
}

func (opt *ListOptions) Where(col_name string, col_value interface{}) {
	opt.WhereColumn = col_name
	opt.WhereValue = col_value
}

func (opt *ListOptions) WhereGroup(where_group *WhereGroup) {
	opt.WhereGroups = append(opt.WhereGroups, where_group)
}

func (opt *ListOptions) Search(searchCols []string, searchFor string) {
	opt.SearchColumns = searchCols
	opt.SearchFor = searchFor
}

func (opt *ListOptions) SetTableName(name string) {
	opt.Table = name
}

func (opt *ListOptions) ShouldAddFirstWhereGroupOperator() bool {
	has_where_before := (opt.WhereColumn != "" && opt.WhereValue != nil) ||
		(len(opt.SearchColumns) != 0 && opt.SearchFor != "")
	return has_where_before && len(opt.WhereGroups) > 0
}

func (opt *ListOptions) ShouldAddFirstWhereGroupWhere() bool {
	has_where_before := (opt.WhereColumn != "" && opt.WhereValue != nil) ||
		(len(opt.SearchColumns) != 0 && opt.SearchFor != "")
	return !has_where_before && len(opt.WhereGroups) > 0
}

func (opt ListOptions) BuildQuery() (q string, args []interface{}) {
	q = fmt.Sprintf("SELECT * From %s", opt.Table)

	if opt.WhereColumn != "" && opt.WhereValue != nil {
		q += fmt.Sprintf(" WHERE %s = :whereValue", opt.WhereColumn)
		args = append(args, sql.Named("whereValue", opt.WhereValue))
	} else if len(opt.SearchColumns) != 0 && opt.SearchFor != "" {
		whereClaudes := []string{}
		for _, col := range opt.SearchColumns {
			whereClaudes = append(whereClaudes, fmt.Sprintf("%s Like :search", col))
		}
		q += fmt.Sprintf(" WHERE %s", strings.Join(whereClaudes, " Or "))
		args = append(args, sql.Named("search", "%"+opt.SearchFor+"%"))
	}

	if opt.ShouldAddFirstWhereGroupWhere() {
		q += " WHERE"

	}
	if opt.ShouldAddFirstWhereGroupOperator() {
		q += fmt.Sprintf(" %s", strings.ToUpper(opt.WhereGroups[0].Operator))
	}

	for i, where_group := range opt.WhereGroups {
		query, where_args := where_group.BuildQuery()
		if i != 0 {
			q += fmt.Sprintf(" %s", opt.WhereGroups[i].Operator)
		}
		q += " " + query
		args = append(args, where_args...)
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
