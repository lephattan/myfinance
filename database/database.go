package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myfinance/env"
	"reflect"
	"slices"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
)

type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	MigrateDriver() (string, database.Driver, error)
	Select(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
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
