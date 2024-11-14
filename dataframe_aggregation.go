package grizzly

func (df *DataFrame) GenericCalculation(operation func(series Series) float64) DataFrame {
	var result []Series
	var value float64
	var newSeries Series
	for _, series := range df.Columns {
		if series.DataType != "float" {
			newSeries = Series{
				Name:     series.Name,
				Float:    []float64{},
				String:   []string{""},
				DataType: "string",
			}
			result = append(result, newSeries)
			continue
		}
		value = operation(series)
		newSeries = Series{
			Name:     series.Name,
			Float:    []float64{value},
			String:   make([]string, 0),
			DataType: "float",
		}
		result = append(result, newSeries)
	}
	return DataFrame{result}
}

func (df *DataFrame) CountEmpty() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetSum()
	})
}

func (df *DataFrame) GetMax() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetMax()
	})
}

func (df *DataFrame) GetMin() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetMin()
	})
}

func (df *DataFrame) GetMean() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetMean()
	})
}

func (df *DataFrame) GetMedian() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetMedian()
	})
}

func (df *DataFrame) GetProduct() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetProduct()
	})
}

func (df *DataFrame) GetSum() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetSum()
	})
}

func (df *DataFrame) GetVariance() DataFrame {
	return df.GenericCalculation(func(series Series) float64 {
		return series.GetVariance()
	})
}

func (df *DataFrame) CountWord(columnName string, word string) int {
	series := df.GetColumnByName(columnName)
	return series.CountWord(word)
}
