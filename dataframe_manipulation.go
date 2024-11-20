package grizzly

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

func (df *DataFrame) FilterFloat(seriesName string, condition func(float64) bool) {
	// Verify if series exists
	series := df.GetColumnByName(seriesName)
	if series.Name == "" {
		fmt.Println("Column not found")
	} else if series.DataType != "float" {
		fmt.Println("Not a float")
	} else {
		indexes := series.FilterFloatSeries(condition)
		for i := range df.Columns {
			df.Columns[i].RemoveIndexes(indexes)
		}
	}
	return
}

func (df *DataFrame) FilterString(columnName string, condition func(string) bool) {
	// Verify if series exists
	series := df.GetColumnByName(columnName)
	indexes := series.FilterStringSeries(condition)
	for i := range df.Columns {
		df.Columns[i].RemoveIndexes(indexes)
	}
	return
}

func (df *DataFrame) ApplyFloat(columnName string, operation func(float64) float64) {
	// Retrieve the series
	series := df.GetColumnByName(columnName)

	// Get the length of the data
	numElements := len(series.Float)
	if numElements == 0 {
		return
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Process chunks in parallel
	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numElements {
			end = numElements
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				series.Float[i] = operation(series.Float[i])
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return
}

func (df *DataFrame) ApplyString(columnName string, operation func(string) string) {
	// Retrieve the series
	series := df.GetColumnByName(columnName)

	// Get the length of the data
	numElements := len(series.String)
	if numElements == 0 {
		return
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Process chunks in parallel
	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numElements {
			end = numElements
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				series.String[i] = operation(series.String[i])
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return
}

func (df *DataFrame) ReplaceWholeWord(column, old, new string) {
	name := df.GetColumnByName(column)
	name.ReplaceWholeWord(old, new)
	return
}

func (df *DataFrame) Replace(column, old, new string) {
	name := df.GetColumnByName(column)
	name.Replace(old, new)
}

func (df *DataFrame) DropByIndex(index ...int) {
	var newSeries []Series
	for i, series := range df.Columns {
		if !ArrayContainsInteger(index, i) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
	return
}

func (df *DataFrame) DropByName(name ...string) {
	var newSeries []Series
	for _, series := range df.Columns {
		if !ArrayContainsString(name, series.Name) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
	return
}

func (df *DataFrame) ConvertStringToFloat(names ...string) {
	for _, name := range names {
		series := df.GetColumnByName(name)
		series.ConvertStringToFloat()
	}
	return
}

func (df *DataFrame) ConvertFloatToString(names ...string) {
	for _, name := range names {
		series := df.GetColumnByName(name)
		series.ConvertFloatToString()
	}
	return
}

func (df *DataFrame) ConvertStringToFloatIndex(indexes ...int) {
	for _, index := range indexes {
		series := df.GetColumnByIndex(index)
		series.ConvertStringToFloat()
	}
	return
}

func (df *DataFrame) ConvertFloatToStringIndex(indexes ...int) {
	for _, index := range indexes {
		series := df.GetColumnByIndex(index)
		series.ConvertFloatToString()
	}
	return
}

func (df *DataFrame) SplitColumn(columnName, delimiter string, newColumnNames []string) {
	column := df.GetColumnByName(columnName)

	if column.DataType == "float" {
		panic("Just for string columns")
	}
	if len(newColumnNames) == 0 {
		panic("No new column names provided")
	}

	numElements := column.GetLength()
	numGoroutines := runtime.NumCPU() // Use number of available CPUs for parallelism
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	// Create slices to hold the new column values
	splitValues := make([][]string, len(newColumnNames))
	for i := range splitValues {
		splitValues[i] = make([]string, numElements)
	}

	var wg sync.WaitGroup

	// Split the work among goroutines
	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numElements {
			end = numElements
		}
		if start >= end {
			break // No more elements to process
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				parts := strings.Split(column.String[i], delimiter)
				for j := 0; j < len(newColumnNames); j++ {
					if j < len(parts) {
						splitValues[j][i] = parts[j]
					} else {
						splitValues[j][i] = "" // Fill missing values with an empty string
					}
				}
			}
		}(start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Add new columns to the DataFrame
	for j, newName := range newColumnNames {
		newColumn := NewStringSeries(newName, splitValues[j])
		df.AddSeries(newColumn)
	}

	return
}

func (df *DataFrame) JoinColumns(columnName1, columnName2, delimiter, newColumnName string) {
	// Retrieve the columns to be joined
	column1 := df.GetColumnByName(columnName1)

	column2 := df.GetColumnByName(columnName2)

	// Validate that both columns are string columns
	if column1.DataType != "string" || column2.DataType != "string" {
		panic("JoinColumns is only supported for string columns")
	}

	numElements := column1.GetLength()
	joinedValues := make([]string, numElements)

	// Use goroutines to join columns in parallel for large datasets
	numGoroutines := runtime.NumCPU()
	if numGoroutines > numElements {
		numGoroutines = numElements // Limit the number of goroutines to the number of elements
	}
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numElements {
			end = numElements
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				joinedValues[i] = column1.String[i] + delimiter + column2.String[i]
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Add the new joined column to the DataFrame
	newColumn := NewStringSeries(newColumnName, joinedValues)
	df.AddSeries(newColumn)

	return
}

func (df *DataFrame) SliceRows(offset int, length int) {
	// Ensure offset is within bounds
	if offset < 0 || offset >= len(df.Columns[0].Float) {
		panic("offset out of range")
	}

	// Ensure length is within bounds
	if offset+length > len(df.Columns[0].Float) {
		length = len(df.Columns[0].Float) - offset // Adjust length to max available range
	}

	// Create a new DataFrame to hold the sliced data
	newDf := &DataFrame{
		Columns: make([]Series, len(df.Columns)),
	}

	// Iterate over each Series to slice the data
	for i, series := range df.Columns {
		newSeries := Series{
			Name:     series.Name,
			DataType: series.DataType,
		}

		// Slice based on the DataType
		if series.DataType == "float" {
			newSeries.Float = series.Float[offset : offset+length]
		} else if series.DataType == "string" {
			newSeries.String = series.String[offset : offset+length]
		}

		newDf.Columns[i] = newSeries
	}

	return
}

func (df *DataFrame) SliceColumnsByIndex(indexes ...int) {
	for i, _ := range df.Columns {
		if ArrayContainsInteger(indexes, i) {
			df.Columns = append(df.Columns[:i], df.Columns[i+1:]...)
		}
	}
}

func (df *DataFrame) MergeDataFrame(otherDf DataFrame, defaultValue string) {
	names := df.GetColumnNames()
	otherNames := otherDf.GetColumnNames()
	for _, name := range names {
		if ArrayContainsString(otherNames, name) {
			panic("Column already exists")
		}
	}
	for _, column := range otherDf.Columns {
		df.AddSeriesForced(column, defaultValue)
	}
}

func (df *DataFrame) Concatenate(otherDf DataFrame, defaultValue string) {
	var newColumn Series
	otherNames := otherDf.GetColumnNames()
	names := df.GetColumnNames()
	for _, name := range otherNames {
		if !ArrayContainsString(names, name) {
			newColumn = NewStringSeries(name, []string{})
			df.AddSeriesForced(newColumn, defaultValue)
		}
	}
	names = append(names, otherNames...)
	var series *Series
	var otherSeries *Series
	for _, name := range names {
		series = df.GetColumnByName(name)
		otherSeries = df.GetColumnByName(name)
		if series.DataType == "float" && otherSeries.DataType == "string" {
			df.ConvertFloatToString(name)
		} else if series.DataType == "string" && otherSeries.DataType == "float" {
			df.ConvertStringToFloat(name)
		}

		if series.DataType == "float" {
			series.Float = append(series.Float, otherSeries.Float...)
		} else if series.DataType == "string" {
			series.String = append(series.String, otherSeries.String...)
		}
	}
}

func (df *DataFrame) DuplicateColumn(names ...string) {
	var ptr *Series
	var series Series
	for _, name := range names {
		ptr = df.GetColumnByName(name)
		series = *ptr
		df.AddSeries(series)
	}
	return
}
