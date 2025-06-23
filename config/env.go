package config

import (
	"os"
)

func getString(key, def string) string {
	res := os.Getenv(key)
	if res == "" {
		return def
	}

	return res
}
