package grizzly

import (
	"math"
	"runtime"
	"sync"
)

func (series *Series) FillNaN(newValue float64) {
	if series.DataType == "string" {
		return
	}

	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if math.IsNaN(series.Float[j]) {
					series.Float[j] = newValue
				}
			}
		}(start, end)
	}
	wg.Wait() // Wait for all goroutines to complete
	return
}

func (series *Series) DropNaN() {
	if series.DataType == "string" {
		return // No-op for string data
	}

	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Channel to collect non-NaN values
	resultChan := make(chan []float64, numGoroutines)

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			nonNaN := make([]float64, 0, end-start) // Temporary slice for non-NaN values
			for j := start; j < end; j++ {
				if !math.IsNaN(series.Float[j]) {
					nonNaN = append(nonNaN, series.Float[j])
				}
			}
			resultChan <- nonNaN
		}(start, end)
	}

	// Close the channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from all goroutines
	result := make([]float64, 0, length)
	for nonNaN := range resultChan {
		result = append(result, nonNaN...)
	}

	// Update the Series with the filtered values
	series.Float = result
}
