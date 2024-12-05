package grizzly

import (
	"fmt"
	"runtime"
	"sync"
)

func GrizzlyToMatrix(df DataFrame) ([][]float64, error) {
	// Get the number of rows and columns in the DataFrame
	numRows := df.GetLength()
	numCols := len(df.Columns)

	// Initialize the matrix
	matrix := make([][]float64, numRows)
	for i := 0; i < numRows; i++ {
		matrix[i] = make([]float64, numCols)
	}

	// Determine the number of CPU cores for concurrency
	numGoroutines := runtime.NumCPU()

	// Error channel to capture the first error
	errChan := make(chan error, 1) // Buffered channel to avoid blocking

	// Use a WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Partition columns into chunks to match the number of CPU cores
	chunkSize := (numCols + numGoroutines - 1) / numGoroutines

	// Process chunks in parallel
	for core := 0; core < numGoroutines; core++ {
		start := core * chunkSize
		end := start + chunkSize
		if end > numCols {
			end = numCols
		}

		if start >= end {
			break // No more work to distribute
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for colIndex := start; colIndex < end; colIndex++ {
				column := df.Columns[colIndex]

				// Ensure column is of type float
				if column.DataType != "float" {
					// Non-blocking send of the error (only send the first error)
					select {
					case errChan <- fmt.Errorf("column %q is not of type float", column.Name):
					default:
					}
					return
				}

				// Populate the matrix for this chunk
				for rowIndex := 0; rowIndex < numRows; rowIndex++ {
					matrix[rowIndex][colIndex] = column.Float[rowIndex]
				}
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check if any error occurred during processing
	if err, ok := <-errChan; ok {
		return nil, err
	}

	return matrix, nil
}
