package grizzly

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

func (df *DataFrame) GetColumnByName(name string) Series {
	for _, series := range df.Columns {
		if series.Name == name {
			return series
		}
	}
	panic(fmt.Sprintf("%s not found", name))
}

func (df *DataFrame) GetColumnByIndex(index int) Series {
	return df.Columns[index]
}

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
}

func (df *DataFrame) FilterString(seriesName string, condition func(string) bool) {
	// Verify if series exists
	series := df.GetColumnByName(seriesName)
	if series.Name == "" {
		fmt.Println("Column not found")
	} else if series.DataType != "string" {
		fmt.Println("Not a string")
	} else {
		indexes := series.FilterStringSeries(condition)
		for i := range df.Columns {
			df.Columns[i].RemoveIndexes(indexes)
		}
	}
}

func (df *DataFrame) ReplaceWholeWord(column, old, new string) {
	name := df.GetColumnByName(column)
	name.ReplaceWholeWord(old, new)
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
}

func (df *DataFrame) DropByName(name ...string) {
	var newSeries []Series
	for _, series := range df.Columns {
		if !ArrayContainsString(name, series.Name) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
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
