package commons

import "syscall"

func EnvString(key, fallback string) string {
	value, ok := syscall.Getenv(key)
	if ok {
		return value
	}
	return fallback
}
