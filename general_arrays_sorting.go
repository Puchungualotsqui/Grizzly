package grizzly

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

func ParallelSortFloat(arr []float64) []float64 {
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
		sortedResult = mergeFloat(sortedResult, sortedChunk)
	}

	return sortedResult
}

func mergeFloat(left, right []float64) []float64 {
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

// parallelSortString performs a parallel sort on an array of strings.
func parallelSortString(arr []string) []string {
	n := len(arr)
	if n <= 1 {
		return arr // Already sorted
	}

	numCPUs := runtime.NumCPU()
	chunkSize := int(math.Ceil(float64(n) / float64(numCPUs)))

	// Channel to collect sorted chunks
	chunks := make(chan []string, numCPUs)

	// Use a WaitGroup to synchronize Goroutines
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
		go func(subarray []string) {
			defer wg.Done()
			sort.Strings(subarray) // Sort the chunk
			chunks <- subarray     // Send the sorted chunk to the channel
		}(arr[start:end])
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Collect and merge sorted chunks
	sortedResult := make([]string, 0, n)
	for sortedChunk := range chunks {
		sortedResult = mergeString(sortedResult, sortedChunk)
	}

	return sortedResult
}

// mergeString merges two sorted string slices into one sorted slice.
func mergeString(left, right []string) []string {
	result := make([]string, 0, len(left)+len(right))
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
