package grizzly

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func (df *DataFrame) ExportToCSV(filePath string) error {
	// Open file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Determine the number of rows and columns
	numRows := df.GetLength()
	numCols := len(df.Columns)

	// Write header (column names)
	header := make([]string, numCols)
	for i, col := range df.Columns {
		header[i] = col.Name
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numRows + numGoroutines - 1) / numGoroutines // Ceiling division

	// Channel for errors
	errorChan := make(chan error, numGoroutines)

	// Worker function
	worker := func(start, end int, errorChan chan<- error) {
		for i := start; i < end; i++ {
			row := make([]string, numCols)
			for j, col := range df.Columns {
				switch col.DataType {
				case "float":
					if i < len(col.Float) {
						row[j] = strconv.FormatFloat(col.Float[i], 'f', -1, 64)
					} else {
						row[j] = ""
					}
				case "string":
					if i < len(col.String) {
						row[j] = col.String[i]
					} else {
						row[j] = ""
					}
				}
			}
			if err := writer.Write(row); err != nil {
				errorChan <- fmt.Errorf("failed to write row: %w", err)
				return
			}
		}
	}

	// Launch workers
	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := (g + 1) * chunkSize
		if end > numRows {
			end = numRows
		}
		go worker(start, end, errorChan)
	}

	// Wait for workers to complete
	for g := 0; g < numGoroutines; g++ {
		if err := <-errorChan; err != nil {
			return err
		}
	}

	return nil
}

func (df *DataFrame) ExportToCSVSimple(filePath string) error {
	// Open file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Determine the number of rows and columns
	maxRows := df.GetLength()
	maxColumns := len(df.Columns)

	// Write header (column names)
	header := make([]string, maxColumns)
	for i, col := range df.Columns {
		header[i] = col.Name
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write rows
	for i := 0; i < maxRows; i++ {
		row := make([]string, maxColumns)
		for j, col := range df.Columns {
			switch col.DataType {
			case "float":
				if i < len(col.Float) {
					row[j] = strconv.FormatFloat(col.Float[i], 'f', -1, 64)
				} else {
					row[j] = ""
				}
			case "string":
				if i < len(col.String) {
					row[j] = col.String[i]
				} else {
					row[j] = ""
				}
			}
		}
		err := writer.Write(row)
		if err != nil {
			return fmt.Errorf("error writing row %d: %w", i, err)
		}

	}

	return nil
}
