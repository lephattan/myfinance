package env

import (
	"os"
	"strings"
)

const (
	PROD Env = "production"
	DEV  Env = "development"
)

type Env string

func (e Env) String() string {
	return string(e)
}

// Read environment variable `key`
// return `def` value if env not found
// panic if `key` is found but not in supported envs
func ReadEnv(key string, def Env) Env {
	v := GetEnv(key, def.String())
	if v == "" {
		return def
	}

	env := Env(strings.ToLower(v))
	switch env {
	case PROD, DEV: // allowed
	default:
		panic("Unexpected app env " + v)
	}
	return env
}

// Get env var of `key`
// return `def` if env `key` has no value
func GetEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
