package helper

import (
	"database/sql"
	"reflect"
)

// Dynamic Scan rows to dest struct
func ModelScan(dest interface{}, rows *sql.Rows) error {
	var v reflect.Value
	v = reflect.ValueOf(dest).Elem()
	pointers := []any{}
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Addr().Interface()
		pointers = append(pointers, value)
	}
	return rows.Scan(pointers...)
}

// // Dynamic Scan rows to slice of model struct
func ModelListScan(dest interface{}, rows *sql.Rows) error {
	dest_type := reflect.TypeOf(dest).Elem()
	dest_elem := dest_type.Elem().Elem()
	cp := reflect.MakeSlice(dest_type, 0, 0)

	for rows.Next() {
		t := reflect.New(dest_elem).Interface() // yay
		if err := ModelScan(t, rows); err != nil {
			return err
		}
		cp = reflect.Append(cp, reflect.ValueOf(t))
	}

	if cp.Len() == 0 {
		return sql.ErrNoRows
	}

	reflect.Indirect(reflect.ValueOf(dest)).Set(cp)
	return rows.Err()
}
