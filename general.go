package grizzly

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
