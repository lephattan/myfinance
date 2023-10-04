package helper

import (
	"database/sql"
	"reflect"
)

// Dynamic Scan rows to dest struct
func ModelScan(dest interface{}, rows *sql.Rows) error {
	v := reflect.ValueOf(dest).Elem()
	pointers := []any{}
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Addr().Interface()
		pointers = append(pointers, value)
	}
	return rows.Scan(pointers...)
}
