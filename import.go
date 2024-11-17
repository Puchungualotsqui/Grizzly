package grizzly

import (
	"encoding/csv"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// ImportCSV reads a CSV file and creates a DataFrame with dynamic parallelism based on the number of CPUs
func ImportCSV(filepath string) DataFrame {
	file, err := os.Open(filepath)
	if err != nil {
		panic("File was not found")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading CSV file")
	}

	if len(records) == 0 {
		return DataFrame{}
	}

	headers := records[0]
	numCols := len(headers)
	numRows := len(records) - 1 // Exclude header row
	columns := make([]Series, numCols)

	// Initialize Series for each header
	for i, header := range headers {
		columns[i] = Series{
			Name:     header,
			DataType: "string",                   // Default type
			String:   make([]string, 0, numRows), // Preallocate memory
			Float:    make([]float64, 0, numRows),
		}
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numCols + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numCols {
			end = numCols
		}

		go func(start, end int) {
			defer wg.Done()
			for colIndex := start; colIndex < end; colIndex++ {
				for _, row := range records[1:] {
					value := row[colIndex]
					// Attempt to parse as float
					if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
						columns[colIndex].Float = append(columns[colIndex].Float, floatValue)
						columns[colIndex].DataType = "float"
					} else {
						columns[colIndex].String = append(columns[colIndex].String, value)
					}
				}
			}
		}(start, end)
	}
	wg.Wait()

	return DataFrame{Columns: columns}
}

func ImportCSVOld(filepath string) DataFrame {
	file, err := os.Open(filepath)
	if err != nil {
		panic("File was not found")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if len(records) == 0 {
		return DataFrame{}
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

	return DataFrame{Columns: columns}
}
