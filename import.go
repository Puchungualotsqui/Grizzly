package grizzly

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func ImportCSV(filepath string) (DataFrame, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return DataFrame{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return DataFrame{}, fmt.Errorf("failed to read CSV file: %v", err)
	}

	size := len(records)
	if size == 0 {
		return DataFrame{}, nil
	}

	headers := records[0]
	rows := records[1:]
	numCols := len(headers)
	numRows := len(rows)
	columns := make([]Series, numCols)

	// Initialize Series for each header
	for i, header := range headers {
		columns[i] = Series{
			Name:     header,
			DataType: "string",                // Default type
			String:   make([]string, numRows), // Preallocate with length
		}
	}

	var result DataFrame

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numRows + numGoroutines - 1) / numGoroutines

	// Create local trackers for each goroutine
	localTrackers := make([][]bool, numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		localTrackers[g] = make([]bool, numCols)
		for i := range localTrackers[g] {
			localTrackers[g][i] = true
		}
	}

	var wg sync.WaitGroup

	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if start >= numRows {
			break // Ensure we don't start beyond the slice length
		}
		if end > numRows {
			end = numRows // Adjust end index for the last chunk
		}

		wg.Add(1)
		go func(start, end, g int) {
			defer wg.Done()
			localTracker := localTrackers[g] // Each goroutine uses its own tracker
			for j := start; j < end; j++ {
				for i := range columns {
					stringValue := ""
					if i < len(rows[j]) {
						stringValue = rows[j][i]
					}
					if stringValue == "" {
						stringValue = "NaN"
					}
					if localTracker[i] {
						if _, err := strconv.ParseFloat(stringValue, 64); err != nil {
							localTracker[i] = false
						}
					}
					columns[i].String[j] = stringValue
				}
			}
		}(start, end, g)
	}

	wg.Wait()

	// Merge local trackers into a global tracker
	globalTracker := make([]bool, numCols)
	for i := 0; i < numCols; i++ {
		globalTracker[i] = true
		for g := 0; g < numGoroutines; g++ {
			if !localTrackers[g][i] {
				globalTracker[i] = false
				break
			}
		}
	}

	// Convert columns to float if they are numeric
	for i := range columns {
		if globalTracker[i] {
			columns[i].ConvertStringToFloat()
		}
	}

	result.Columns = columns
	result.FixShape()

	return result, nil
}

func ImportCSVOld(filepath string) (DataFrame, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return DataFrame{}, fmt.Errorf("file was not found")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if len(records) == 0 {
		return DataFrame{}, nil
	}

	headers := records[0]
	columns := make([]Series, len(headers))

	// Initialize Series for each header
	for i, header := range headers {
		columns[i] = Series{Name: header, DataType: "string", String: []string{}}
	}

	// Populate Series with data
	for _, row := range records[1:] {
		for i, value := range row {
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				columns[i].Float = append(columns[i].Float, floatValue)
				columns[i].DataType = "float"
			} else {
				columns[i].String = append(columns[i].String, value)
			}
		}
	}

	return DataFrame{Columns: columns}, nil
}
