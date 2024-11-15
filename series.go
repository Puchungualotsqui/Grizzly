package grizzly

import (
	"fmt"
	"strconv"
)

type Series struct {
	Name     string
	Float    []float64
	String   []string
	DataType string
}

func NewStringSeries(name string, String []string) Series {
	return Series{
		Name:     name,
		Float:    make([]float64, 0),
		String:   String,
		DataType: "string",
	}
}

func NewFloatSeries(name string, Float []float64) Series {
	return Series{
		Name:     name,
		Float:    Float,
		String:   make([]string, 0),
		DataType: "float",
	}
}

func (series *Series) ResizeSeries(targetLength int, defaultValue string) {
	series.ConvertFloatToString()
	series.String = ArrayResizeString(series.String, targetLength, defaultValue)
}

func (series *Series) Print(max int) {
	max = MinInt(max, series.GetLength())
	if series.DataType == "float" {
		for i := 0; i < max; i++ {
			fmt.Println(i, series.Float[i])
		}
	} else {
		for i := 0; i < max; i++ {
			fmt.Println(i, series.String[i])
		}
	}
}

func (series *Series) GetValueFloat(index int) float64 {
	return series.Float[index]
}

func (series *Series) GetValueString(index int) string {
	return series.String[index]
}

func (series *Series) GetValueAsString(index int) string {
	if series.DataType == "float" {
		return strconv.FormatFloat(series.Float[index], 'f', -1, 64)
	} else {
		return series.GetValueString(index)
	}
}
