package grizzly

import (
	"math"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

func ArrayCountWord(data []string, word string) float64 {
	chain := ArrayStringBase(0, data, func(info string, result float64) float64 {
		if info == word {
			result++
		}
		return result
	})
	var result float64

	for val := range chain {
		result += val
	}
	return result
}

func ParallelSort(arr []float64) []float64 {
	n := len(arr)
	if n <= 1 {
		return arr
	}

	numCPUs := runtime.NumCPU()
	chunkSize := int(math.Ceil(float64(n) / float64(numCPUs)))

	// Channel to collect sorted chunks
	chunks := make(chan []float64, numCPUs)

	// Use a WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	for i := 0; i < numCPUs; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= n {
			break // Avoid processing if start is beyond the array length
		}
		if end > n {
			end = n // Ensure we don't go out of bounds
		}

		wg.Add(1)

		// Sort each chunk in a separate Goroutine
		go func(subarray []float64) {
			defer wg.Done()
			sort.Float64s(subarray) // Sort the chunk
			chunks <- subarray      // Send it to the channel
		}(arr[start:end])
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Collect and merge sorted chunks
	sortedResult := make([]float64, 0, n)
	for sortedChunk := range chunks {
		sortedResult = Merge(sortedResult, sortedChunk)
	}

	return sortedResult
}

func Merge(left, right []float64) []float64 {
	result := make([]float64, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	// Append any remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func ArrayCountStringDuplicates(elements []string) map[string]int {
	numCPU := runtime.NumCPU()
	length := len(elements)
	if length == 0 {
		return nil // Handle empty input
	}

	// Calculate chunk size for splitting work
	chunkSize := (length + numCPU - 1) / numCPU
	resultChan := make(chan map[string]int, numCPU)
	var wg sync.WaitGroup

	// Worker function to count duplicates in a subset of elements
	countDuplicates := func(subset []string) {
		defer wg.Done()
		localCounts := make(map[string]int)
		for _, element := range subset {
			localCounts[element]++
		}
		resultChan <- localCounts // Send local result to the channel
	}

	// Start goroutines to process chunks of the array
	for i := 0; i < numCPU && i*chunkSize < length; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go countDuplicates(elements[start:end])
	}

	// Close the channel once all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Merge results from all goroutines
	combinedCounts := make(map[string]int)
	for localCounts := range resultChan {
		for key, count := range localCounts {
			combinedCounts[key] += count
		}
	}

	// Filter results to retain only elements with counts greater than 1
	finalResult := make(map[string]int)
	for key, count := range combinedCounts {
		if count > 1 {
			finalResult[key] = count
		}
	}

	return finalResult
}

func ArrayCountFloatDuplicates(elements []float64) map[float64]int {
	numCPU := runtime.NumCPU()
	length := len(elements)
	if length == 0 {
		return nil // Handle empty input
	}

	// Calculate chunk size for splitting work
	chunkSize := (length + numCPU - 1) / numCPU
	resultChan := make(chan map[float64]int, numCPU)
	var wg sync.WaitGroup

	// Worker function to count duplicates in a subset of elements
	countDuplicates := func(subset []float64) {
		defer wg.Done()
		localCounts := make(map[float64]int)
		for _, element := range subset {
			localCounts[element]++
		}
		resultChan <- localCounts // Send local result to the channel
	}

	// Start goroutines to process chunks of the array
	for i := 0; i < numCPU && i*chunkSize < length; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go countDuplicates(elements[start:end])
	}

	// Close the channel once all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Merge results from all goroutines
	combinedCounts := make(map[float64]int)
	for localCounts := range resultChan {
		for key, count := range localCounts {
			combinedCounts[key] += count
		}
	}

	// Filter results to retain only elements with counts greater than 1
	finalResult := make(map[float64]int)
	for key, count := range combinedCounts {
		if count > 1 {
			finalResult[key] = count
		}
	}

	return finalResult
}

func ArrayContainsInteger(arr []int, target int) bool {
	for _, value := range arr {
		if value == target {
			return true // Element found
		}
	}
	return false // Element not found
}

func ArrayContainsString(arr []string, target string) bool {
	for _, value := range arr {
		if value == target {
			return true // Element found
		}
	}
	return false // Element not found
}

func ArrayGetNonFloatValues(input []string) []string {
	var nonConvertible []string

	for _, str := range input {
		_, err := strconv.ParseFloat(str, 64)
		if err != nil {
			nonConvertible = append(nonConvertible, str) // Collect non-convertible elements
		}
	}

	return nonConvertible
}

func ArrayResizeString(input []string, targetLength int, defaultValue string) []string {
	for len(input) < targetLength {
		input = append(input, defaultValue)
	}
	return input
}
