package grizzly

import (
	"encoding/csv"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// ImportCSV reads a CSV file and creates a DataFrame with consistent column lengths
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

	size := len(records)
	if size == 0 {
		return DataFrame{}
	}

	headers := records[0]
	rows := records[1:]
	numCols := len(headers)
	numRows := size - 1 // Exclude header row
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
	var result DataFrame

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numRows + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= numRows {
			break // Ensure we don't start beyond the slice length
		}
		if end > numRows {
			end = numRows // Adjust end index for the last chunk
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				for i, _ := range columns {
					columns[i].String[j] = rows[j][i]
				}
			}
		}(start, end)
	}

	// Closing channel after all goroutines finish
	go func() {
		wg.Wait()
	}()

	for i, _ := range columns {
		columns[i].ConvertStringToFloat()
	}

	result.Columns = columns
	result.FixShape("")

	return result
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
