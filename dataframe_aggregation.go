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

func (df *DataFrame) CountWord(word string) DataFrame {
	var columns []Series
	var count []float64
	var result DataFrame
	for _, series := range df.Columns {
		count[0] = float64(series.CountWord(word))
		columns = append(columns, NewFloatSeries(series.Name, count))
	}
	result.Columns = columns
	result.FixShape("")
	return result
}

func (df *DataFrame) GetNonFloatValues() DataFrame {
	var result []Series
	var tempSeries Series
	var tempString []string
	var resultDataframe DataFrame
	for _, column := range df.Columns {
		tempString = column.GetNonFloatValues()
		if len(tempString) != 0 {
			tempSeries = NewStringSeries(column.Name, tempString)
			result = append(result, tempSeries)
		}
	}
	resultDataframe = DataFrame{result}
	resultDataframe.FixShape("")
	return resultDataframe
}

func (df *DataFrame) GetUniqueValues() DataFrame {
	var result DataFrame
	var temp Series
	var tempString []string
	var tempFloat []float64
	for _, column := range df.Columns {
		if column.DataType == "string" {
			tempString = ArrayUniqueValuesString(column.String)
			temp = NewStringSeries(column.Name, tempString)
			result.Columns = append(result.Columns, temp)
			tempString = nil
		} else {
			tempFloat = ArrayUniqueValuesFloat(column.Float)
			temp = NewFloatSeries(column.Name, tempFloat)
			result.Columns = append(result.Columns, temp)
			tempFloat = nil
		}
	}
	result.FixShape("")
	return result
}
