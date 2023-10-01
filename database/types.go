package database

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
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

// Custom int64 type for storing date in unix format
type UnixDate struct {
	int64
}

// Format UnixDate to string with given format
func (u *UnixDate) Format(format string) string {
	t := time.Unix(u.int64, 0)
	return t.Format(format)
}

// Format UnixDate as "02-01-2006"
func (u *UnixDate) String() string {
	return u.Format("02-01-2006")
}

// int64 value of UnixDate
func (u *UnixDate) Init64() int64 {
	return u.int64
}

// Scan given value into UnixDate
func (u *UnixDate) Scan(value interface{}) error {
	if i, ok := value.(int64); ok {
		u.int64 = int64(i)
	}
	return nil
}

// Value for sql value conversion
func (u UnixDate) Value() (driver.Value, error) {
	return u.int64, nil
}

// Convert string format of 2006-01-02 to reflect.Value of UnixDate
func UnixDateConverter(value string) reflect.Value {
	v, err := time.Parse("2006-01-02", value)
	if err != nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(UnixDate{v.Unix()})
}

var UnixDateParser = fiber.ParserType{
	Customtype: UnixDate{},
	Converter:  UnixDateConverter,
}
