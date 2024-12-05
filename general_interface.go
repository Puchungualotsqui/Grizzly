package grizzly

import (
	"fmt"
	"strconv"
)

func interfaceConvertToString(value any) (string, error) {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%.2f", v), nil
	case int, int32, int64:
		return fmt.Sprintf("%v", v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func interfaceConvertToFloat(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, fmt.Errorf("unsupported type: %T", value)
}
