package env

import (
	"time"
)

func GetString(key string) string {
	return immutable.GetString(key)
}

func GetStringSlice(key string) []string {
	return immutable.GetStringSlice(key)
}

func GetInt(key string) int {
	return immutable.GetInt(key)
}

func GetIntSlice(key string) []int {
	return immutable.GetIntSlice(key)
}

func GetUint(key string) uint {
	return immutable.GetUint(key)
}

func GetDuration(key string) time.Duration {
	return immutable.GetDuration(key)
}

func GetFloat64(key string) float64 {
	return immutable.GetFloat64(key)
}

type sentinelDefaultType interface {
	~string | []string | ~int | []int | ~bool | ~uint | ~float64 | time.Duration
}

// GetWithDefault get the value from the environment if it exists, otherwise return the default value.
func GetWithDefault[T sentinelDefaultType](key string, defaultValue T) T {
	if immutable.IsSet(key) {
		switch any(defaultValue).(type) {
		case string:
			return any(immutable.GetString(key)).(T)
		case []string:
			return any(immutable.GetStringSlice(key)).(T)
		case int:
			return any(immutable.GetInt(key)).(T)
		case uint:
			return any(immutable.GetUint(key)).(T)
		case []int:
			return any(immutable.GetIntSlice(key)).(T)
		case bool:
			return any(immutable.GetBool(key)).(T)
		case float64:
			return any(immutable.GetFloat64(key)).(T)
		case time.Duration:
			return any(immutable.GetDuration(key)).(T)
		}
	}

	return defaultValue
}
