package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myfinance/env"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
)

type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	MigrateDriver() (string, database.Driver, error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)
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

// Return DB based on "APP_ENV" environment variable
func GetDB() DB {
	app_env := env.ReadEnv("APP_ENV", "production")
	db := NewDB(app_env)
	return db
}

// Return DB besed on given Env
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

var ErrUnprocessable = errors.New("invalid model")

func GetAffectedRows(result sql.Result) int {
	if result == nil {
		return 0
	}
	n, _ := result.RowsAffected()
	return int(n)
}

func GenerateInsertStatement(m Record) (stmt string, args []interface{}, err error) {

	tablename := m.TableName()
	tag_name := "db"
	stmt = fmt.Sprintf("Insert Into %s ", tablename)

	type_m := reflect.TypeOf(m)
	m_values := reflect.ValueOf(m)

	args = []any{}
	cols := []string{}

	for i := 0; i < type_m.NumField(); i++ {
		value := m_values.Field(i).Interface()
		struct_field := type_m.Field(i)
		db_tags := strings.Split(string(struct_field.Tag.Get(tag_name)), ",")
		if slices.Contains(db_tags, "omitinsert") {
			continue
		}
		db_tag := db_tags[0]
		if db_tag == "" {
			continue
		}

		cols = append(cols, db_tag)
		args = append(args, value)
	}
	if len(cols) > 0 {
		placeholders := make([]string, len(args))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		stmt = stmt + fmt.Sprintf("(%s) Values (%s)", strings.Join(cols, ", "), strings.Join(placeholders, ","))
	}

	return
}
