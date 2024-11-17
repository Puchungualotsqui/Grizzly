package grizzly

import (
	"runtime"
	"sync"
)

func ArrayFloatBase(initValue float64, data []float64, operation func(info float64, result float64) float64) chan float64 {
	length := len(data)
	if length == 0 {
		// Handle empty data case by returning a closed channel immediately
		emptyChan := make(chan float64)
		close(emptyChan)
		return emptyChan
	}
	numGoroutines := runtime.NumCPU()
	if numGoroutines > length {
		numGoroutines = length // Avoid creating more goroutines than necessary
	}
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	resultChan := make(chan float64, numGoroutines)

	// Function to calculate the sum of a chunk
	worker := func(start, end int) {
		defer wg.Done()
		result := initValue
		// Always starts from second value to calculate Mean Correctly
		for i := start; i < end; i++ {
			result = operation(data[i], result)
		}
		resultChan <- result
	}

	// Launch goroutines to process chunks
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func ArrayStringBase(initValue float64, data []string, operation func(info string, result float64) float64) chan float64 {
	length := len(data)
	if length == 0 {
		// Handle empty data case by returning a closed channel immediately
		emptyChan := make(chan float64)
		close(emptyChan)
		return emptyChan
	}
	numGoroutines := runtime.NumCPU()
	if numGoroutines > length {
		numGoroutines = length // Avoid creating more goroutines than necessary
	}
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	resultChan := make(chan float64, numGoroutines)

	// Function to calculate the sum of a chunk
	worker := func(start, end int) {
		defer wg.Done()
		result := initValue
		// Always starts from second value to calculate Mean Correctly
		for i := start; i < end; i++ {
			result = operation(data[i], result)
		}
		resultChan <- result
	}

	// Launch goroutines to process chunks
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func IsNameRepeated(seriesArray []Series, targetName string) bool {
	for _, s := range seriesArray {
		if s.Name == targetName {
			return true
		}
	}
	return false
}
