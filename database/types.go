package database

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

// Nullable Int64 in sql
// Reference: https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
type NullInt64 sql.NullInt64

// NullInt64 into value
// validate value if nil or not
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}

	return nil
}

func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}
