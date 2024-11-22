package grizzly

import (
	"math"
)

func ArrayMean(data []float64) float64 {
	chain := ArrayFloatBase(0, data, func(info float64, result float64) float64 {
		result = result + info
		return result
	})
	var result float64

	for val := range chain {
		result += val
	}
	result /= float64(len(data))
	return result
}

func ArrayProduct(data []float64) float64 {
	chain := ArrayFloatBase(1, data, func(info float64, result float64) float64 {
		result = result * info
		return result
	})
	result := 1.0

	for val := range chain {
		result *= val
	}
	return result
}

func ArraySum(data []float64) float64 {
	chain := ArrayFloatBase(0, data, func(info float64, result float64) float64 {
		result = result + info
		return result
	})
	var result float64

	for val := range chain {
		result += val
	}
	return result
}

func ArrayVariance(data []float64) float64 {
	mean := ArrayMean(data)
	chain := ArrayFloatBase(0, data, func(info float64, result float64) float64 {
		diff := info - mean
		return result + diff*diff // Accumulate the squared difference
	})

	var sumOfSquaredDiffs float64
	for val := range chain {
		sumOfSquaredDiffs += val
	}

	// Step 3: Calculate the variance (sum of squared differences divided by the number of elements)
	return sumOfSquaredDiffs / float64(len(data))
}

func ArrayMin(data []float64) float64 {
	maxChan := ArrayFloatBase(math.MaxFloat64, data, func(info float64, result float64) float64 {
		if info < result {
			result = info
		}
		return result
	})

	minVal := <-maxChan // Initialize minVal with the first value received from the channel
	for val := range maxChan {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

func ArrayMax(data []float64) float64 {
	maxChan := ArrayFloatBase(math.MaxFloat64*-1, data, func(info float64, result float64) float64 {
		if info > result {
			result = info
		}
		return result
	})

	minVal := <-maxChan // Initialize minVal with the first value received from the channel
	for val := range maxChan {
		if val > minVal {
			minVal = val
		}
	}
	return minVal
}

func ArrayMedian(nums []float64) float64 {
	nums = ParallelSortFloat(nums)
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
