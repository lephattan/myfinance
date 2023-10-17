package helper

import (
	"strconv"
	"time"
)

// Parse string into uint64
// Return 0 on empty string
// Return 0 on error
func ParseUint64(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		i = 0
	}
	return i, err
}

// Format unix timestamp into string with given format
func UnixTimeFmt(unixT int64, format string) string {
	t := time.Unix(unixT, 0)
	return t.Format(format)
}

type number interface {
	int | uint | int64 | uint64
}

func Devide[T number](a, b T) T {
	var s T
	s = a / b
	return s
}

func Minus[T number](a, b T) T {
	return a - b
}
