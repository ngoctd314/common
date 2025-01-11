package env

import (
	"fmt"
	"strings"
	"time"
)

func required(key string) string {
	return fmt.Sprintf("%s is required", strings.ToUpper(strings.ReplaceAll(key, ".", "_")))
}

func MustString(key string) string {
	if immutable.IsSet(key) {
		return immutable.GetString(key)
	}
	panic(required(key))
}

func MustStringSlice(key string) []string {
	if immutable.IsSet(key) {
		return immutable.GetStringSlice(key)
	}
	panic(required(key))
}

func MustInt(key string) int {
	if immutable.IsSet(key) {
		return immutable.GetInt(key)
	}
	panic(required(key))
}

func MustIntSlice(key string) []int {
	if immutable.IsSet(key) {
		return immutable.GetIntSlice(key)
	}
	panic(required(key))
}

func MustUint(key string) uint {
	if immutable.IsSet(key) {
		return immutable.GetUint(key)
	}
	panic(required(key))
}

func MustDuration(key string) time.Duration {
	if immutable.IsSet(key) {
		return immutable.GetDuration(key)
	}
	panic(required(key))
}

func MustFloat64(key string) float64 {
	if immutable.IsSet(key) {
		return immutable.GetFloat64(key)
	}
	panic(required(key))
}
