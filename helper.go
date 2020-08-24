package ubot

import "reflect"

func withDefault(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	rValue := reflect.ValueOf(value)
	if !rValue.IsValid() || rValue.IsZero() {
		return defaultValue
	}
	return value
}
