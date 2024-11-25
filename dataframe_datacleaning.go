package grizzly

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"
)

func (df *DataFrame) FillNaN(newValue float64) {
	for i := range df.Columns {
		df.Columns[i].FillNaN(newValue)
	}
}

func (df *DataFrame) DropNaN() {
	for i := range df.Columns {
		df.Columns[i].DropNaN()
	}
}

func (df *DataFrame) RemoveOutliersZScore(identifier any, threshold float64) error {
	var series *Series
	var zScoreSeries Series
	var err error
	series, err = df.GetColumnDynamic(identifier)
	if err != nil {
		return err
	}
	if series.DataType == "string" {
		return fmt.Errorf("%v is an string column. Please select an float column", identifier)
	}
	var mean float64
	var stdDev float64
	zScore := make([]float64, df.GetLength())

	mean = ArrayMean(series.Float)
	stdDev = ArrayVariance(series.Float)
	stdDev = math.Sqrt(stdDev)

	for i, value := range series.Float {
		zScore[i] = (value - mean) / stdDev
	}
	zScoreSeries = NewFloatSeries("__zScore__", zScore)
	err = df.AddSeries(zScoreSeries)
	if err != nil {
		return err
	}

	condition := func(value float64) bool {
		if value >= threshold {
			return true
		}
		return false
	}
	err = df.FilterFloat("__zScore__", condition)
	if err != nil {
		return err
	}
	df.DropByName("__zScore__")
	return nil
}

func (df *DataFrame) RemoveOutliersIQR(identifier any) error {
	var err error
	var series *Series

	series, err = df.GetColumnDynamic(identifier)
	if err != nil {
		return err
	}
	if series.DataType == "string" {
		return fmt.Errorf("%v is an string column. Please select an float column", identifier)
	}
	err = df.Sort(identifier)
	if err != nil {
		return err
	}

	q1 := ArrayCalculatePercentile(series.Float, 25)
	q3 := ArrayCalculatePercentile(series.Float, 75)
	iqr := q3 - q1

	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	operation := func(value float64) bool {
		if value <= lowerBound || value >= upperBound {
			return true
		}
		return false
	}
	err = df.FilterFloat(identifier, operation)
	if err != nil {
		return err
	}
	return nil
}

func (df *DataFrame) RemoveDuplicates() {
	if len(df.Columns) == 0 {
		return // No data
	}

	// Determine the number of rows
	rowCount := len(df.Columns[0].Float)
	if df.Columns[0].DataType == "string" {
		rowCount = len(df.Columns[0].String)
	}

	numGoroutines := runtime.NumCPU() // Number of concurrent workers
	chunkSize := (rowCount + numGoroutines - 1) / numGoroutines

	globalSeen := sync.Map{}  // Concurrent map to track unique rows
	uniqueIndices := []int{}  // Indices of unique rows
	indicesMu := sync.Mutex{} // Mutex to protect uniqueIndices
	var wg sync.WaitGroup     // WaitGroup to manage goroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > rowCount {
			end = rowCount
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			localUniqueIndices := []int{}

			for idx := start; idx < end; idx++ {
				row := make([]interface{}, len(df.Columns))

				// Collect row data across all columns
				for j, col := range df.Columns {
					if col.DataType == "float" {
						row[j] = col.Float[idx]
					} else if col.DataType == "string" {
						row[j] = col.String[idx]
					}
				}

				// Serialize the row to JSON for unique identification
				rowKey, err := json.Marshal(row)
				if err != nil {
					fmt.Println("Error serializing row:", err)
					continue
				}

				// Check if the row is globally unique
				if _, loaded := globalSeen.LoadOrStore(string(rowKey), true); !loaded {
					localUniqueIndices = append(localUniqueIndices, idx)
				}
			}

			// Merge local unique indices into global uniqueIndices
			indicesMu.Lock()
			uniqueIndices = append(uniqueIndices, localUniqueIndices...)
			indicesMu.Unlock()
		}(start, end)
	}

	wg.Wait() // Wait for all goroutines to complete

	// Sort uniqueIndices to maintain row order (optional but recommended)
	sort.Ints(uniqueIndices)

	// Create new slices for unique data
	for j := range df.Columns {
		if df.Columns[j].DataType == "float" {
			newFloatData := make([]float64, len(uniqueIndices))
			for k, idx := range uniqueIndices {
				newFloatData[k] = df.Columns[j].Float[idx]
			}
			df.Columns[j].Float = newFloatData
		} else if df.Columns[j].DataType == "string" {
			newStringData := make([]string, len(uniqueIndices))
			for k, idx := range uniqueIndices {
				newStringData[k] = df.Columns[j].String[idx]
			}
			df.Columns[j].String = newStringData
		}
	}
}
