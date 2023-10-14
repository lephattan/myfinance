package database

import (
	"database/sql/driver"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Implementation of SQL nullable with generic
// Ref: https://github.com/golang/go/issues/60370
type Nullable[T any] struct {
	Actual T
	Valid  bool
}

func convertAssign[T any](out *T, in any) error {
	*out = in.(T)
	return nil
}

// Scan value into Nullable
func (n *Nullable[T]) Scan(value any) error {
	n.Valid = false
	if value == nil {
		var v T
		n.Actual = v
		return nil
	}
	err := convertAssign[T](&n.Actual, value)
	if err == nil {
		n.Valid = true
	}
	return err
}

// Return value for sql driver to exec in query
func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Actual, nil
}

// Convert input string to Nullable[int64]
func NullableInt64Converter(value string) reflect.Value {
	null_value := reflect.ValueOf(Nullable[int64]{
		Actual: 0,
		Valid:  false,
	})

	if value == "" {
		return null_value
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return null_value
	}

	return reflect.ValueOf(Nullable[int64]{
		Actual: v,
		Valid:  true,
	})
}

// Custom parser for Nullable[int64]
var NullableInt64Parser = fiber.ParserType{
	Customtype: Nullable[int64]{},
	Converter:  NullableInt64Converter,
}
