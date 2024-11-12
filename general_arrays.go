package grizzly

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

func ArrayCountEmpty(words []string) int {
	numGoroutines := runtime.NumCPU() // Use the number of available CPU cores
	length := len(words)
	if length == 0 {
		return 0
	}

	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	countChan := make(chan int, numGoroutines) // Buffered channel to collect counts

	// Launch goroutines to count empty strings in chunks
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			count := 0
			for j := start; j < end; j++ {
				if words[j] == "" {
					count++
				}
			}
			countChan <- count // Send the count to the channel
		}(start, end)
	}

	// Close the channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(countChan)
	}()

	// Collect the counts from all goroutines
	totalCount := 0
	for count := range countChan {
		totalCount += count
	}

	return totalCount

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
	chunkSize := (len(elements) + numCPU - 1) / numCPU // Calculate chunk size for splitting work
	duplicateCounts := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Worker function to count duplicates in a subset of elements
	countDuplicates := func(subset []string) {
		defer wg.Done()
		localCounts := make(map[string]int)
		for _, element := range subset {
			localCounts[element]++
		}
		// Safely merge localCounts into the global map
		mu.Lock()
		for key, count := range localCounts {
			duplicateCounts[key] += count
		}
		mu.Unlock()
	}

	// Start goroutines to process chunks of the array
	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= len(elements) {
			break // Prevent out-of-bound access if start index is beyond the length
		}
		if end > len(elements) {
			end = len(elements) // Adjust end index to the length of the array
		}
		wg.Add(1)
		go countDuplicates(elements[start:end])
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Filter results to retain only elements with counts greater than 1
	result := make(map[string]int)
	for key, count := range duplicateCounts {
		if count > 1 {
			result[key] = count
		}
	}

	return result
}

func ArrayCountFloatDuplicates(elements []float64) map[float64]int {
	numCPU := runtime.NumCPU()
	chunkSize := (len(elements) + numCPU - 1) / numCPU // Calculate chunk size for splitting work
	duplicateCounts := make(map[float64]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Worker function to count duplicates in a subset of elements
	countDuplicates := func(subset []float64) {
		defer wg.Done()
		localCounts := make(map[float64]int)
		for _, element := range subset {
			localCounts[element]++
		}
		// Safely merge localCounts into the global map
		mu.Lock()
		for key, count := range localCounts {
			duplicateCounts[key] += count
		}
		mu.Unlock()
	}

	// Start goroutines to process chunks of the array
	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= len(elements) {
			break // Prevent out-of-bound access if start index is beyond the length
		}
		if end > len(elements) {
			end = len(elements) // Adjust end index to the length of the array
		}
		wg.Add(1)
		go countDuplicates(elements[start:end])
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Filter results to retain only elements with counts greater than 1
	result := make(map[float64]int)
	for key, count := range duplicateCounts {
		if count > 1 {
			result[key] = count
		}
	}

	return result
}
