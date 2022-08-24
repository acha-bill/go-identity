package utils

import "os"

// GetEnv returns the env variable defined by key or returns the default if specified.
func GetEnv(key string, def ...string) string {
	if len(def) == 0 {
		return os.Getenv(key)
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return def[0]
	}
	return val
}
