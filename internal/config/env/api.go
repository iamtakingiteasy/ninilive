// Package env provides environment-based configuration
package env

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/iamtakingiteasy/ninilive/internal/config"
)

// NewLoader returns env config loader
func NewLoader() config.Loader {
	return &loader{}
}

// String resolves string env variable
func String(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return strings.TrimSpace(v)
	}

	return def
}

// StringArray resolves string array env variable
func StringArray(name string, sep *regexp.Regexp, def []string) []string {
	if v, ok := os.LookupEnv(name); ok {
		return sep.Split(strings.TrimSpace(v), -1)
	}

	return def
}

// Uint64 resolves uint64 env variable
func Uint64(name string, def uint64) uint64 {
	if raw, ok := os.LookupEnv(name); ok {
		if v, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64); err == nil {
			return v
		}
	}

	return def
}

// Int64 resolves int64 env variable
func Int64(name string, def int64) int64 {
	if raw, ok := os.LookupEnv(name); ok {
		if v, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64); err == nil {
			return v
		}
	}

	return def
}

// Bool resolves bool env variable
func Bool(name string, def bool) bool {
	if raw, ok := os.LookupEnv(name); ok {
		if v, err := strconv.ParseBool(strings.TrimSpace(raw)); err == nil {
			return v
		}
	}

	return def
}
