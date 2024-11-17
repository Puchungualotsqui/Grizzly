package grizzly

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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
	if len(df.Columns) == 0 {
		df.Columns = append(df.Columns, series)
	}
	if series.GetLength() != df.GetLength() {
		fmt.Println("ERROR: Cannot add a series with different length")
	} else if IsNameRepeated(df.Columns, series.Name) {
		fmt.Println("ERROR: Cannot add a series with repeat names")
	}
}

func (df *DataFrame) FixShape(defaultValue string) {
	var size int
	for _, series := range df.Columns {
		size = MaxInt(size, series.GetLength())
	}
	for i, _ := range df.Columns {
		df.Columns[i].ResizeSeries(size, defaultValue)
	}
}

func (df *DataFrame) AddSeriesForced(series Series, defaultValue string) {
	df.Columns = append(df.Columns, series)
	df.FixShape(defaultValue)
}

func (df *DataFrame) Print(min int, max int) {
	// Ensure max does not exceed the length of the DataFrame
	max = MinInt(df.GetLength(), max)

	// Create a tabwriter for better column alignment
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Add "Index" as the first header
	names := df.GetColumnNames()
	headers := append([]string{"Index"}, names...)
	fmt.Fprintln(writer, strings.Join(headers, "\t"))

	// Print rows of data
	for i := min; i < max; i++ {
		var output []string
		// Add the row index as the first element
		output = append(output, strconv.Itoa(i))
		// Add column values for this row
		for _, series := range df.Columns {
			output = append(output, series.GetValueAsString(i))
		}
		// Join and print the row with the tabwriter
		fmt.Fprintln(writer, strings.Join(output, "\t"))
	}

	// Flush the writer to ensure output is printed
	writer.Flush()
}

func (df *DataFrame) PrintHead(max int) {
	df.Print(0, max)
}

func (df *DataFrame) PrintTail(max int) {
	df.Print(max, df.GetLength())
}

func (df *DataFrame) PrintOld(max int) {
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
