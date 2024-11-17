package grizzly

import (
	"encoding/csv"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// ImportCSV reads a CSV file using streaming processing and creates a DataFrame with optimized performance
func ImportCSV(filepath string) DataFrame {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Error reading CSV headers: %v", err)
	}
	numCols := len(headers)

	// Initialize Series for each column with empty data
	columns := make([]Series, numCols)
	for i, header := range headers {
		columns[i] = Series{
			Name:     header,
			DataType: "string",
			String:   []string{},
			Float:    []float64{},
		}
	}

	// Determine the number of goroutines for parallel processing
	numGoroutines := runtime.NumCPU()
	rowChannel := make(chan []string, numGoroutines*2) // Buffered channel for rows
	var wg sync.WaitGroup

	// Launch goroutines to process rows
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowChannel {
				for i, value := range row {
					// Parse value and append to the appropriate column
					if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
						columns[i].Float = append(columns[i].Float, floatValue)
						columns[i].DataType = "float"
					} else {
						columns[i].String = append(columns[i].String, value)
					}
				}
			}
		}()
	}

	// Stream rows from the CSV file
	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			log.Printf("Error reading row: %v", err)
			continue
		}
		rowChannel <- row // Send row for processing
	}
	close(rowChannel)
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
