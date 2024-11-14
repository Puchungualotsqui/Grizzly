package grizzly

func (series *Series) CountWord(word string) int {
	if series.DataType == "float" {
		return 0
	} else {
		return ArrayCountWord(series.String, word)
	}
}

func (series *Series) GetMax() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetMax requires a non-empty array")
	}
	return ArrayMax(series.Float)
}

func (series *Series) GetMin() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetMin requires a non-empty array")
	}
	return ArrayMin(series.Float)
}

func (series *Series) GetMean() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetMean requires a non-empty array")
	}
	return ArrayMean(series.Float)
}

func (series *Series) GetMedian() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetMedian requires a non-empty array")
	}
	return ArrayMedian(series.Float)
}

func (series *Series) GetProduct() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetProduct requires a non-empty array")
	}
	return ArrayProduct(series.Float)
}

func (series *Series) GetSum() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetSum requires a non-empty array")
	}
	return ArraySum(series.Float)
}

func (series *Series) GetVariance() float64 {
	if series.DataType == "string" {
		return 0
	} else if series.GetLength() == 0 {
		panic("GetSum requires a non-empty array")
	}
	return ArrayVariance(series.Float)
}
