package grizzly

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func (series *Series) RemoveIndexes(indexes []int) {
	if series.DataType == "float" {
		filteredFloats := make([]float64, len(indexes))
		for i, idx := range indexes {
			filteredFloats[i] = series.Float[idx]
		}
		series.Float = filteredFloats
	} else {
		filteredStrings := make([]string, len(indexes))
		for i, idx := range indexes {
			filteredStrings[i] = series.String[idx]
		}
		series.String = filteredStrings
	}
}

func (series *Series) ConvertStringToFloat() {
	if series.DataType == "float" {
		return
	}
	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	length := len(series.String)
	floatArray := make([]float64, length)
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	// Calculate chunk size
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Launch multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}

		wg.Add(1)

		// Process the chunk in a goroutine
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				// Stop if an error has already occurred
				if firstErr != nil {
					return
				}

				// Try to convert the string to a float
				val, err := strconv.ParseFloat(series.String[j], 64)
				if err != nil {
					once.Do(func() {
						firstErr = err
					})
					return // Stop processing this goroutine
				}
				floatArray[j] = val
			}
		}(start, end)
	}
	wg.Wait()

	// Check if an error occurred during conversion
	if firstErr != nil {
		fmt.Println("Processing stopped due to error: ", firstErr)
	} else {
		series.Float = floatArray
		series.String = nil // Clear the string slice
		series.DataType = "float"
	}
	return
}

func (series *Series) ConvertFloatToString() {
	if series.DataType == "string" {
		return
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	length := len(series.Float)
	stringArray := make([]string, length)
	var wg sync.WaitGroup

	// Calculate chunk size for splitting the work among goroutines
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Launch multiple goroutines
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
				// Convert float to string and store in the string array
				stringArray[j] = strconv.FormatFloat(series.Float[j], 'f', -1, 64)
			}
		}(start, end)
	}
	wg.Wait()

	// Update the series with the new string data
	series.String = stringArray
	series.Float = nil // Clear the float slice
	series.DataType = "string"
	return
}

func (series *Series) ReplaceWholeWord(old, new string) {
	if series.DataType == "float" || series.GetLength() == 0 {
		return
	}

	numGoroutines := runtime.NumCPU()
	length := series.GetLength()
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	// Compile the regular expression once
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(old))
	re := regexp.MustCompile(pattern)

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
				// No mutex needed; each goroutine works on separate slice elements
				series.String[j] = re.ReplaceAllString(series.String[j], new)
			}
		}(start, end)
	}
	wg.Wait() // Wait for all goroutines to complete
	return
}

func (series *Series) Replace(old, new string) {
	if series.DataType == "float" || series.GetLength() == 0 {
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
				series.String[j] = strings.ReplaceAll(series.String[j], old, new)
			}
		}(start, end)
	}
	wg.Wait() // Wait for all goroutines to complete
	return
}
