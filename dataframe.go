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

func (df *DataFrame) CreateFloatColumn(name string, nums []float64) error {
	newSeries := NewFloatSeries(name, nums)
	err := df.AddSeries(newSeries)
	if err != nil {
		// Propagate the error if adding the series fails
		return fmt.Errorf("failed to create float column %q: %w", name, err)
	}
	// Successfully added the column
	return nil
}

func (df *DataFrame) CreateStringColumn(name string, words []string) error {
	newSeries := NewStringSeries(name, words)
	err := df.AddSeries(newSeries)
	if err != nil {
		// Propagate the error if adding the series fails
		return fmt.Errorf("failed to create string column %q: %w", name, err)
	}
	// Successfully added the column
	return nil
}

func (df *DataFrame) AddSeries(series Series) error {
	if len(df.Columns) == 0 {
		df.Columns = append(df.Columns, series)
		return nil
	}
	if series.GetLength() != df.GetLength() {
		return fmt.Errorf("cannot add a series with different length: series length %d, dataframe length %d",
			series.GetLength(), df.GetLength())

	} else if isNameRepeated(df.Columns, series.Name) {
		return fmt.Errorf("cannot add a series with repeated name: %s", series.Name)
	}
	df.Columns = append(df.Columns, series)
	return nil
}

func (df *DataFrame) FixShape() {
	var size int
	for _, series := range df.Columns {
		size = maxInt(size, series.GetLength())
	}
	for i := range df.Columns {
		df.Columns[i].ResizeSeries(size, "NaN")
	}
}

func (df *DataFrame) AddSeriesForced(series Series) {
	df.Columns = append(df.Columns, series)
	df.FixShape()
}

func (df *DataFrame) Print(min int, max int) error {
	var err error
	// Ensure max does not exceed the length of the DataFrame
	max = minInt(df.GetLength(), max)

	// Create a tab writer for better column alignment
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Add "Index" as the first header
	names := df.GetColumnNames()
	headers := append([]string{"Index"}, names...)
	_, err = fmt.Fprintln(writer, strings.Join(headers, "\t"))
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Print rows of data
	for i := min; i < max; i++ {
		var output []string
		// Add the row index as the first element
		output = append(output, strconv.Itoa(i))
		// Add column values for this row
		for _, series := range df.Columns {
			output = append(output, series.GetValueAsString(i))
		}
		// Join and print the row with the tab writer
		_, err = fmt.Fprintln(writer, strings.Join(output, "\t"))
		if err != nil {
			return fmt.Errorf("failed to write row %d: %w", i, err)
		}
	}

	// Flush the writer to ensure output is printed
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}

func (df *DataFrame) PrintHead(max int) error {
	if max <= 0 {
		return fmt.Errorf("invalid value for max: %d (must be > 0)", max)
	}

	// Call Print and handle any errors
	if err := df.Print(0, max); err != nil {
		return fmt.Errorf("failed to print head: %w", err)
	}
	return nil
}

func (df *DataFrame) PrintTail(min int) error {
	if min < 0 || min >= df.GetLength() {
		return fmt.Errorf("invalid value for min: %d (must be >= 0 and < dataframe length)", min)
	}

	// Call Print and handle any errors
	if err := df.Print(min, df.GetLength()); err != nil {
		return fmt.Errorf("failed to print tail: %w", err)
	}
	return nil
}
