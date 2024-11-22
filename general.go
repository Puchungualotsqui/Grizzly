package grizzly

import (
	"strconv"
	"strings"
)

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func IsNameRepeated(seriesArray []Series, targetName string) bool {
	for _, s := range seriesArray {
		if s.Name == targetName {
			return true
		}
	}
	return false
}

func LengthOfFloat(value float64) (int, int) {
	// Convert the float to a string
	str := strconv.FormatFloat(value, 'f', -1, 64) // Convert with full precision
	parts := strings.Split(str, ".")

	// Count digits before the decimal point
	beforeDecimal := len(parts[0])

	// Count digits after the decimal point if it exists
	afterDecimal := 0
	if len(parts) > 1 {
		afterDecimal = len(parts[1])
	}

	return beforeDecimal, afterDecimal
}
