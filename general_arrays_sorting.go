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
		sortedResult = MergeFloat(sortedResult, sortedChunk)
	}

	return sortedResult
}

func MergeFloat(left, right []float64) []float64 {
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

type Element struct {
	Value float64
	Index int
}

func ParallelSortFloatMap(arr []float64) ([]float64, []int) {
	n := len(arr)
	if n <= 1 {
		return arr, make([]int, n)
	}

	// Create an array of Element structs to track values and original indices
	elements := make([]Element, n)
	for i, v := range arr {
		elements[i] = Element{Value: v, Index: i}
	}

	numCPUs := runtime.NumCPU()
	chunkSize := int(math.Ceil(float64(n) / float64(numCPUs)))

	// Channel to collect sorted chunks
	chunks := make(chan []Element, numCPUs)

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
		go func(subarray []Element) {
			defer wg.Done()
			sort.Slice(subarray, func(i, j int) bool {
				return subarray[i].Value < subarray[j].Value
			}) // Sort the chunk
			chunks <- subarray // Send it to the channel
		}(elements[start:end])
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Collect and merge sorted chunks
	sortedResult := make([]Element, 0, n)
	for sortedChunk := range chunks {
		sortedResult = MergeFloatElements(sortedResult, sortedChunk)
	}

	// Extract the index mapping for the sorted order
	indexMapping := make([]int, n)
	for i, elem := range sortedResult {
		indexMapping[i] = elem.Index
	}

	// Extract the sorted values
	sortedValues := make([]float64, n)
	for i, elem := range sortedResult {
		sortedValues[i] = elem.Value
	}

	return sortedValues, indexMapping
}

func MergeFloatElements(left, right []Element) []Element {
	result := make([]Element, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i].Value < right[j].Value {
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

type StringElement struct {
	Value string
	Index int
}

func ParallelSortStringsMap(arr []string) ([]string, []int) {
	n := len(arr)
	if n <= 1 {
		return arr, make([]int, n)
	}

	// Create an array of StringElement structs to track values and original indices
	elements := make([]StringElement, n)
	for i, v := range arr {
		elements[i] = StringElement{Value: v, Index: i}
	}

	numCPUs := runtime.NumCPU()
	chunkSize := int(math.Ceil(float64(n) / float64(numCPUs)))

	// Channel to collect sorted chunks
	chunks := make(chan []StringElement, numCPUs)

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
		go func(subarray []StringElement) {
			defer wg.Done()
			sort.Slice(subarray, func(i, j int) bool {
				return subarray[i].Value < subarray[j].Value
			}) // Sort the chunk
			chunks <- subarray // Send it to the channel
		}(elements[start:end])
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Collect and merge sorted chunks
	sortedResult := make([]StringElement, 0, n)
	for sortedChunk := range chunks {
		sortedResult = MergeStringElements(sortedResult, sortedChunk)
	}

	// Extract the index mapping for the sorted order
	indexMapping := make([]int, n)
	for i, elem := range sortedResult {
		indexMapping[i] = elem.Index
	}

	// Extract the sorted values
	sortedValues := make([]string, n)
	for i, elem := range sortedResult {
		sortedValues[i] = elem.Value
	}

	return sortedValues, indexMapping
}

func MergeStringElements(left, right []StringElement) []StringElement {
	result := make([]StringElement, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i].Value < right[j].Value {
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
