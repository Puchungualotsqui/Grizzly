package grizzly

import (
	"fmt"
	"strconv"
	"strings"
)

type DataFrame struct {
	Columns []Series
}

func CreateDataFrame(series ...Series) DataFrame {
	return DataFrame{
		Columns: series,
	}
}

func (df *DataFrame) CreateFloatColumn(name string, nums []float64) {
	new_series := NewFloatSeries(name, nums)
	df.AddSeries(new_series)
}

func (df *DataFrame) CreateStringColumn(name string, words []string) {
	new_series := NewStringSeries(name, words)
	df.AddSeries(new_series)
}

func (df *DataFrame) AddSeries(series Series) {
	if series.GetLength() != df.GetLength() {
		fmt.Println("ERROR: Cannot add a series with different length")
	} else if IsNameRepeated(df.Columns, series.Name) {
		fmt.Println("ERROR: Cannot add a series with repeat names")
	} else {
		df.Columns = append(df.Columns, series)
	}
}

func (df *DataFrame) Print(max int) {
	max = MinInt(df.GetLength(), max)
	var output []string
	names := df.GetColumnNames()
	for _, name := range names {
		output = append(output, name)
	}
	fmt.Println(strings.Join(output, "  "))
	output = output[:0]
	for i := 0; i < max; i++ {
		output = append(output, strconv.Itoa(i))
		for _, series := range df.Columns {
			output = append(output, series.GetValueAsString(i))
		}
		fmt.Println(strings.Join(output, "  "))
		output = output[:0]
	}
}
