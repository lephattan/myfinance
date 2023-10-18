package helper

import (
	"log"
	"math"
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

// Convert interface into int
// return 0 if value is not in uint64, uint, int64, int
func ToInt(value interface{}) int {
	var i int
	switch v := value.(type) {
	case int:
		i = v
	case int64:
		i = int(v)
	case uint:
		i = int(v)
	case uint64:
		i = int(v)
	default:
		i = 0
	}
	return i
}

// Add two number, support uint64, uint, int64, int
// add 0 if cannot convert interface value into int
func Add(a, b interface{}) int {
	return ToInt(a) + ToInt(b)
}

// Create a slice of int64 stating from 0
// Sequence(LAST)
// Sequence(START, LAST)
// Sequence(START, STEP, LAST)
func Sequence(args ...uint64) []uint64 {
	log.Printf("Args: %v", args)
	var start, last, step, length uint64
	step = 1
	start = 0
	if len(args) == 1 {
		last = args[0]
	} else if len(args) == 2 {
		start = args[0]
		last = args[1]
	} else if len(args) == 3 {
		start = args[0]
		step = args[1]
		last = args[2]
	}

	length = uint64(math.Floor(float64(last-start) / float64(step)))

	seq := make([]uint64, length)
	for inx := range seq {
		seq[inx] = start + step*uint64(inx)
	}
	return seq

}
