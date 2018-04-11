package gateway

import (
	"strconv"
	"fmt"
)

func Int642(val interface{}) (int64, error) {
	switch t := val.(type) {
	case int:
		// TODO
	case string:
		// TODO
	default:
		return 0, fmt.Errorf("gateway: type '%s' is not int or string", t)
	}
	return 0, nil
}

// String just returns the given string.
// It is just for compatibility to other types.
func String(val string) (string, error) {
	return val, nil
}

// Bool converts the given string representation of a boolean value into bool.
func Bool(val string) (bool, error) {
	return strconv.ParseBool(val)
}

// Float64 converts the given string representation into representation of a floating point number into float64.
func Float64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

// Int64 converts the given string representation of an integer into int64.
func Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

// Int32 converts the given string representation of an integer into int32.
func Int32(val string) (int32, error) {
	i, err := strconv.ParseInt(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// Int16 converts the given string representation of an integer into int16.
func Int16(val string) (int16, error) {
	i, err := strconv.ParseInt(val, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

// Int8 converts the given string representation of an integer into int8.
func Int8(val string) (int8, error) {
	i, err := strconv.ParseInt(val, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}
