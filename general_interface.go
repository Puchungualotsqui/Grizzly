package grizzly

import "fmt"

func InterfaceConvertToString(value any) (string, error) {
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
