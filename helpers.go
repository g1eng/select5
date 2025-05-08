package select5

import (
	"fmt"
)

// GetS extracts a string value from an any type.
// Returns an error if the value cannot be converted to a string.
func GetS(v any) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	case *string:
		if v.(*string) == nil {
			return "", nil
		}
		return *v.(*string), nil
	default:
		return "", fmt.Errorf("invalid type: %T", v)
	}
}

// GetI extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int.
func GetI(v any) (int, error) {
	switch v.(type) {
	case int:
		return v.(int), nil
	case *int:
		if v.(*int) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*int), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetI8 extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int8.
func GetI8(v any) (int8, error) {
	switch v.(type) {
	case int8:
		return v.(int8), nil
	case *int8:
		if v.(*int8) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*int8), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetI16 extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int16.
func GetI16(v any) (int16, error) {
	switch v.(type) {
	case int16:
		return v.(int16), nil
	case *int16:
		if v.(*int16) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*int16), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetI32 extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int32.
func GetI32(v any) (int32, error) {
	switch v.(type) {
	case int32:
		return v.(int32), nil
	case *int32:
		if v.(*int32) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*int32), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetI64 extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int64.
func GetI64(v any) (int64, error) {
	switch v.(type) {
	case int64:
		return v.(int64), nil
	case *int64:
		if v.(*int64) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*int64), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetF32 extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int32.
func GetF32(v any) (float32, error) {
	switch v.(type) {
	case float32:
		return v.(float32), nil
	case *float32:
		if v.(*float32) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*float32), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetF64 extracts an int value from an any type.
// Returns an error if the value cannot be converted to an int64.
func GetF64(v any) (float64, error) {
	switch v.(type) {
	case float64:
		return v.(float64), nil
	case *float64:
		if v.(*float64) == nil {
			return 0, fmt.Errorf("no data: %T", v)
		}
		return *v.(*float64), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", v)
	}
}

// GetB extracts a bool value from an any type.
// Returns an error if the value cannot be converted to a bool.
func GetB(v any) (bool, error) {
	switch v.(type) {
	case bool:
		return v.(bool), nil
	case *bool:
		if v.(*bool) == nil {
			return false, nil
		}
		return *v.(*bool), nil
	default:
		return false, fmt.Errorf("invalid type: %T", v)
	}
}

// GetV is a generic extract function for interface values.
// Returns the value as string if it can be converted, or an error otherwise.
func GetV(v any) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	case []byte:
		return string(v.([]byte)), nil
	case int:
		return fmt.Sprintf("%d", v.(int)), nil
	case int8:
		return fmt.Sprintf("%d", v.(int8)), nil
	case int16:
		return fmt.Sprintf("%d", v.(int16)), nil
	case int32:
		return fmt.Sprintf("%d", v.(int32)), nil
	case int64:
		return fmt.Sprintf("%d", v.(int64)), nil
	case float32:
		return fmt.Sprintf("%f", v.(float32)), nil
	case float64:
		return fmt.Sprintf("%f", v.(float64)), nil
	case bool:
		if v.(bool) {
			return "✓", nil
		} else {
			return "", nil
		}
	default:
		return "", fmt.Errorf("type %T not supported", v)
	}
}

// GetVP is a generic extract function for interface values for pointer types.
// Returns the value as string if it can be converted, or an error otherwise.
func GetVP(v any) (string, error) {
	switch v.(type) {
	case string:
		p, err := GetS(v)
		if err != nil {
			return "", err
		}
		return p, nil
	case *string:
		p, err := GetS(v)
		if err != nil {
			return "", err
		}
		return p, nil
	case []byte:
		return string(v.([]byte)), nil
	case int:
		i, err := GetI(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case *int:
		i, err := GetI(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case int8:
		i, err := GetI8(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case *int8:
		i, err := GetI8(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case int16:
		i, err := GetI16(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case *int16:
		i, err := GetI16(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case int32:
		i, err := GetI32(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case *int32:
		i, err := GetI32(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case int64:
		i, err := GetI64(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case *int64:
		i, err := GetI64(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case float32:
		f, err := GetF32(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%f", f), nil
	case *float32:
		f, err := GetF32(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%f", f), nil
	case float64:
		f, err := GetF64(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%f", f), nil
	case *float64:
		f, err := GetF64(v)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%f", f), nil
	case bool:
		b, err := GetB(v)
		if err != nil {
			return "", err
		}
		if b {
			return "✓", nil
		} else {
			return "", nil
		}
	case *bool:
		b, err := GetB(v)
		if err != nil {
			return "", err
		}
		if b {
			return "✓", nil
		} else {
			return "", nil
		}
	default:
		return "", fmt.Errorf("type %T not supported", v)
	}
}

func CheckPrimitive(s any) (res byte) {
	switch s.(type) {
	case string:
		res |= IsString
	case int:
		res |= IsInt
	case int8:
		res |= IsInt8
	case int16:
		res |= IsInt16
	case int32:
		res |= IsInt32
	case int64:
		res |= IsInt64
	case float32:
		res |= IsFloat32
	case float64:
		res |= IsFloat64
	case bool:
		res |= IsBool
	case *string:
		res |= IsPointer | IsString
	case *int:
		res |= IsPointer | IsInt
	case *int8:
		res |= IsPointer | IsInt8
	case *int16:
		res |= IsPointer | IsInt16
	case *int32:
		res |= IsPointer | IsInt32
	case *int64:
		res |= IsPointer | IsInt64
	case *float32:
		res |= IsPointer | IsFloat32
	case *float64:
		res |= IsPointer | IsFloat64
	case *bool:
		res |= IsPointer | IsBool
	default:
		res |= IsAny
	}
	return
}
