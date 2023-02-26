package core

import "reflect"

func ConditionalInt(condition bool, onTrue, onFalse int) int {
	if condition {
		return onTrue
	} else {
		return onFalse
	}
}

// IsNil checks if input is nil or not.
//
// It will return "true" if input interface is wrapping a nil pointer variable.
//
// It will return "false" If interface is pointing to a non-pointer value.
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	t := reflect.TypeOf(i)
	isValueType := t.Kind() != reflect.Ptr
	if isValueType {
		return false
	}

	return reflect.ValueOf(i).IsNil()
}
