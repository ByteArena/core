package utils

import (
	"os"
)

func GetenvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetenvOrThrow(name string) string {
	value := os.Getenv(name)

	Assert(value != "", "ERROR: could not find "+name+" in env")

	return value
}
