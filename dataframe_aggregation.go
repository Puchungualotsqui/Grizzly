package grizzly

import (
	"fmt"
)

func (df *DataFrame) GenericCalculation(operation func(series Series) (float64, error)) (DataFrame, error) {
	var result []Series
	var value float64
	var newSeries Series
	var err error

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
		value, err = operation(series)
		if err != nil {
			return DataFrame{}, fmt.Errorf("error executing calculation")
		}
		newSeries = Series{
			Name:     series.Name,
			Float:    []float64{value},
			String:   make([]string, 0),
			DataType: "float",
		}
		result = append(result, newSeries)
	}
	return DataFrame{result}, nil
}

func (df *DataFrame) GetMax() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetMax()
	})
}

func (df *DataFrame) GetMin() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetMin()
	})
}

func (df *DataFrame) GetMean() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetMean()
	})
}

func (df *DataFrame) GetMedian() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetMedian()
	})
}

func (df *DataFrame) GetProduct() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetProduct()
	})
}

func (df *DataFrame) GetSum() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetSum()
	})
}

func (df *DataFrame) GetVariance() (DataFrame, error) {
	return df.GenericCalculation(func(series Series) (float64, error) {
		return series.GetVariance()
	})
}

func (df *DataFrame) CountWord(word string) (DataFrame, error) {
	var count float64
	var result DataFrame
	var err error
	number, numberToo := TryConvertToFloat(word)
	for _, series := range df.Columns {
		if series.DataType == "float" {
			if numberToo == true {
				count = ArrayFloatCountValue(series.Float, number)
			} else {
				count = 0
			}
		} else {
			count = ArrayStringCountWord(series.String, word)
		}
		err = result.CreateFloatColumn(series.Name, []float64{count})
		if err != nil {
			return DataFrame{}, fmt.Errorf("failed to create column %q: %w", series.Name, err)
		}
	}

	return result, err
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
	resultDataframe.FixShape()
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
	result.FixShape()
	return result
}

func (df *DataFrame) CountNaNValues() DataFrame {
	var count float64
	series := make([]Series, len(df.Columns))
	for i, column := range df.Columns {
		series[i].DataType = "float"
		series[i].Name = column.Name
		if column.DataType == "float" {
			count = ArrayFloatCountNaNValue(column.Float)
			series[i].Float = []float64{count}

		} else {
			count = ArrayStringCountWord(column.String, "NaN")
			series[i].Float = []float64{count}
		}
	}
	return DataFrame{series}
}
