// Package shared: for shared functions
package shared

import "os"

func GetEnvString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
