package grizzly

import (
	"runtime"
	"sync"
)

func ArrayMedian(nums []float64) float64 {
	nums = ParallelSort(nums)
	n := len(nums)
	if n == 0 {
		panic("array cannot be empty")
	}

	if n%2 == 1 {
		// Odd length, return the middle element
		return nums[n/2]
	} else {
		// Even length, return the average of the two middle elements
		return (nums[n/2-1] + nums[n/2]) / 2.0
	}
}

func ArrayMean(nums []float64) float64 {
	length := len(nums)

	numGoroutines := runtime.NumCPU()
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	sumChan := make(chan float64, numGoroutines)

	// Worker function to calculate partial sum
	calculatePartialSum := func(start, end int) {
		defer wg.Done()
		localSum := 0.0
		for i := start; i < end; i++ {
			localSum += nums[i]
		}
		sumChan <- localSum
	}

	// Split work among goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > length {
			end = length
		}
		if start < end { // Only start goroutine if there's work to do
			wg.Add(1)
			go calculatePartialSum(start, end)
		}
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(sumChan)
	}()

	// Collect the results and calculate the total sum
	totalSum := 0.0
	for partialSum := range sumChan {
		totalSum += partialSum
	}

	mean := totalSum / float64(length)
	return mean
}

func ArrayMax(nums []float64) float64 {
	length := len(nums)
	numGoroutines := runtime.NumCPU()
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	maxChan := make(chan float64, numGoroutines)

	// Worker function to find the max in a chunk
	findMax := func(start, end int) {
		defer wg.Done()
		localMax := nums[start]
		for i := start + 1; i < end; i++ {
			if nums[i] > localMax {
				localMax = nums[i]
			}
		}
		maxChan <- localMax
	}

	// Split work among goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > length {
			end = length
		}
		if start < end { // Only start goroutine if there's work to do
			wg.Add(1)
			go findMax(start, end)
		}
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(maxChan)
	}()

	// Collect the results and find the overall max
	maxVal := <-maxChan
	for localMax := range maxChan {
		if localMax > maxVal {
			maxVal = localMax
		}
	}

	return maxVal
}

func ArrayMin(nums []float64) float64 {
	length := len(nums)
	numGoroutines := runtime.NumCPU()
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	maxChan := make(chan float64, numGoroutines)

	// Worker function to find the max in a chunk
	findMin := func(start, end int) {
		defer wg.Done()
		localMin := nums[start]
		for i := start + 1; i < end; i++ {
			if nums[i] < localMin {
				localMin = nums[i]
			}
		}
		maxChan <- localMin
	}

	// Split work among goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > length {
			end = length
		}
		if start < end { // Only start goroutine if there's work to do
			wg.Add(1)
			go findMin(start, end)
		}
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(maxChan)
	}()

	// Collect the results and find the overall max
	minVal := <-maxChan
	for localMin := range maxChan {
		if localMin < minVal {
			minVal = localMin
		}
	}

	return minVal
}

func ArrayProduct(nums []float64) float64 {
	n := len(nums)
	if n == 0 {
		return 0 // ArrayProduct of an empty list is defined as 1
	}

	// Determine the number of CPU cores available
	numCPU := runtime.NumCPU()
	chunkSize := (n + numCPU - 1) / numCPU // Calculate chunk size to distribute work
	results := make(chan float64, numCPU)  // Channel to collect partial results
	var wg sync.WaitGroup

	// Function to calculate the product of a chunk
	worker := func(start, end int) {
		defer wg.Done()
		product := 1.0
		for i := start; i < end; i++ {
			product *= nums[i]
		}
		results <- product
	}

	// Launch goroutines to process chunks
	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate the results from all goroutines
	totalProduct := 1.0
	for product := range results {
		totalProduct *= product
	}

	return totalProduct
}

func ArraySum(nums []float64) float64 {
	n := len(nums)
	if n == 0 {
		return 0 // Sum of an empty list is 0
	}

	// Determine the number of CPU cores available
	numCPU := runtime.NumCPU()
	chunkSize := (n + numCPU - 1) / numCPU // Calculate chunk size to distribute work
	results := make(chan float64, numCPU)  // Channel to collect partial results
	var wg sync.WaitGroup

	// Function to calculate the sum of a chunk
	worker := func(start, end int) {
		defer wg.Done()
		sum := 0.0
		for i := start; i < end; i++ {
			sum += nums[i]
		}
		results <- sum
	}

	// Launch goroutines to process chunks
	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate the results from all goroutines
	totalSum := 0.0
	for sum := range results {
		totalSum += sum
	}

	return totalSum
}

func ArrayVariance(nums []float64) float64 {
	n := len(nums)
	if n == 0 {
		return 0 // Return 0 for an empty array
	}

	mean := ArrayMean(nums)

	// Determine the number of CPU cores available
	numCPU := runtime.NumCPU()
	chunkSize := (n + numCPU - 1) / numCPU // Calculate chunk size to distribute work
	results := make(chan float64, numCPU)  // Channel to collect partial results
	var wg sync.WaitGroup

	// Function to calculate the sum of squared differences for a chunk
	worker := func(start, end int) {
		defer wg.Done()
		sumOfSquaredDiffs := 0.0
		for i := start; i < end; i++ {
			diff := nums[i] - mean
			sumOfSquaredDiffs += diff * diff
		}
		results <- sumOfSquaredDiffs
	}

	// Launch goroutines to process chunks
	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate the results from all goroutines
	totalSumOfSquaredDiffs := 0.0
	for sumOfSquaredDiffs := range results {
		totalSumOfSquaredDiffs += sumOfSquaredDiffs
	}

	// Calculate variance
	variance := totalSumOfSquaredDiffs / float64(n)
	return variance
}
