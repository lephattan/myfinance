package helper

import "strconv"

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
