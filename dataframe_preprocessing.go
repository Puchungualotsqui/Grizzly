package grizzly

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

/*#############
#Data Cleaning#
#############*/

func (df *DataFrame) FillNaN(newValue float64, identifiers ...any) error {
	var err error
	if len(identifiers) != 0 {
		var series *Series
		for _, identifier := range identifiers {
			series, err = df.GetColumnDynamic(identifier)
			if err != nil {
				return err
			}
			series.FillNaN(newValue)
		}
		return nil
	}
	for i := range df.Columns {
		df.Columns[i].FillNaN(newValue)
	}
	return nil
}

func (df *DataFrame) DropNaN(identifiers ...any) error {
	var err error
	if len(identifiers) != 0 {
		var series *Series
		for _, identifier := range identifiers {
			series, err = df.GetColumnDynamic(identifier)
			if err != nil {
				return err
			}
			series.DropNaN()
		}
		return nil
	}
	for i := range df.Columns {
		df.Columns[i].DropNaN()
	}
	return nil
}

func (df *DataFrame) RemoveOutliersZScore(identifier any, threshold float64) error {
	var series *Series
	var zScoreSeries Series
	var err error
	series, err = df.GetColumnDynamic(identifier)
	if err != nil {
		return err
	}
	if series.DataType == "string" {
		return fmt.Errorf("%v is an string column. Please select an float column", identifier)
	}
	var mean float64
	var stdDev float64
	zScore := make([]float64, df.GetLength())

	mean = arrayMean(series.Float)
	stdDev = arrayVariance(series.Float)
	stdDev = math.Sqrt(stdDev)

	for i, value := range series.Float {
		zScore[i] = (value - mean) / stdDev
	}
	zScoreSeries = NewFloatSeries("__zScore__", zScore)
	err = df.AddSeries(zScoreSeries)
	if err != nil {
		return err
	}

	condition := func(value float64) bool {
		if value >= threshold {
			return true
		}
		return false
	}
	err = df.FilterFloat("__zScore__", condition)
	if err != nil {
		return err
	}
	df.DropByName("__zScore__")
	return nil
}

func (df *DataFrame) RemoveOutliersIQR(identifier any) error {
	var err error
	var series *Series

	series, err = df.GetColumnDynamic(identifier)
	if err != nil {
		return err
	}
	if series.DataType == "string" {
		return fmt.Errorf("%v is an string column. Please select an float column", identifier)
	}
	err = df.Sort(identifier)
	if err != nil {
		return err
	}

	q1 := arrayCalculatePercentile(series.Float, 25)
	q3 := arrayCalculatePercentile(series.Float, 75)
	iqr := q3 - q1

	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	operation := func(value float64) bool {
		if value <= lowerBound || value >= upperBound {
			return true
		}
		return false
	}
	err = df.FilterFloat(identifier, operation)
	if err != nil {
		return err
	}
	return nil
}

func (df *DataFrame) RemoveDuplicates() {
	if len(df.Columns) == 0 {
		return // No data
	}

	// Determine the number of rows
	rowCount := len(df.Columns[0].Float)
	if df.Columns[0].DataType == "string" {
		rowCount = len(df.Columns[0].String)
	}

	numGoroutines := runtime.NumCPU() // Number of concurrent workers
	chunkSize := (rowCount + numGoroutines - 1) / numGoroutines

	globalSeen := sync.Map{}  // Concurrent map to track unique rows
	uniqueIndices := []int{}  // Indices of unique rows
	indicesMu := sync.Mutex{} // Mutex to protect uniqueIndices
	var wg sync.WaitGroup     // WaitGroup to manage goroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > rowCount {
			end = rowCount
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			localUniqueIndices := []int{}

			for idx := start; idx < end; idx++ {
				row := make([]interface{}, len(df.Columns))

				// Collect row data across all columns
				for j, col := range df.Columns {
					if col.DataType == "float" {
						row[j] = col.Float[idx]
					} else if col.DataType == "string" {
						row[j] = col.String[idx]
					}
				}

				// Serialize the row to JSON for unique identification
				rowKey, err := json.Marshal(row)
				if err != nil {
					fmt.Println("Error serializing row:", err)
					continue
				}

				// Check if the row is globally unique
				if _, loaded := globalSeen.LoadOrStore(string(rowKey), true); !loaded {
					localUniqueIndices = append(localUniqueIndices, idx)
				}
			}

			// Merge local unique indices into global uniqueIndices
			indicesMu.Lock()
			uniqueIndices = append(uniqueIndices, localUniqueIndices...)
			indicesMu.Unlock()
		}(start, end)
	}

	wg.Wait() // Wait for all goroutines to complete

	// Sort uniqueIndices to maintain row order (optional but recommended)
	sort.Ints(uniqueIndices)

	// Create new slices for unique data
	for j := range df.Columns {
		if df.Columns[j].DataType == "float" {
			newFloatData := make([]float64, len(uniqueIndices))
			for k, idx := range uniqueIndices {
				newFloatData[k] = df.Columns[j].Float[idx]
			}
			df.Columns[j].Float = newFloatData
		} else if df.Columns[j].DataType == "string" {
			newStringData := make([]string, len(uniqueIndices))
			for k, idx := range uniqueIndices {
				newStringData[k] = df.Columns[j].String[idx]
			}
			df.Columns[j].String = newStringData
		}
	}
}

/*###############
#Feature Scaling#
###############*/

func (df *DataFrame) Normalize(identifiers ...any) error {
	var minV, maxV float64
	var series *Series
	var err error

	for _, identifier := range identifiers {
		// Retrieve the series (column) by identifier
		series, err = df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("error retrieving column '%v': %w", identifier, err)
		}

		// Ensure the column is numeric
		if series.DataType == "string" {
			return fmt.Errorf("column '%v' is a string column. Please select a float column", identifier)
		}

		if series.GetLength() == 0 {
			return fmt.Errorf("column '%v' is an empty column", identifier)
		}

		// Find minimum and maximum values in the column
		minV = arrayMin(series.Float)
		maxV = arrayMax(series.Float)

		// Handle case where all values are equal
		if minV == maxV {
			series.Float = make([]float64, len(series.Float)) // Set all values to 0
			continue                                          // Move to the next column
		}

		// Normalize the column values
		for i, value := range series.Float {
			series.Float[i] = (value - minV) / (maxV - minV)
		}
	}

	return nil // Successfully normalized
}

func (df *DataFrame) Standardize(identifiers ...any) error {
	var mean, stdDev float64
	var series *Series
	var err error

	for _, identifier := range identifiers {
		series, err = df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("error retrieving column '%v': %w", identifier, err)
		}
		// Ensure the column is numeric
		if series.DataType == "string" {
			return fmt.Errorf("column '%v' is a string column. Please select a float column", identifier)
		}

		if series.GetLength() == 0 {
			return fmt.Errorf("column '%v' is an empty column", identifier)
		}

		mean = arrayMean(series.Float)
		stdDev = arrayVariance(series.Float, mean)
		stdDev = math.Sqrt(stdDev)

		if stdDev == 0 {
			series.Float = make([]float64, len(series.Float))
			continue
		}

		for i, value := range series.Float {
			series.Float[i] = (value - mean) / stdDev
		}
	}
	return nil
}

/*
################################
#Encoding Categorical Variables#
################################
*/
func (df *DataFrame) OneHotEncode(identifiers ...any) error {
	var series *Series
	var err error
	var categories []string
	var lastIndex int
	categoryIndex := make(map[string]int)

	if df.GetLength() == 0 {
		return fmt.Errorf("dataframe is empty")
	}

	numElements := df.GetLength()
	numGoroutines := runtime.NumCPU()
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	for _, identifier := range identifiers {
		series, err = df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("error retrieving column '%v': %w", identifier, err)
		}
		if series.DataType != "string" {
			return fmt.Errorf("column '%v' is not a string column", identifier)
		}

		lastIndex = df.GetNumberOfColumns() - 1

		// Get unique categories from the column
		categories = arrayUniqueValuesString(series.String)

		// Create new float columns for each category with independent slices
		for _, category := range categories {
			// Create an independent empty slice for each column
			emptyColumn := make([]float64, numElements)
			err = df.CreateFloatColumn(category, emptyColumn)
			if err != nil {
				return fmt.Errorf("error creating column '%v': %w", category, err)
			}
			lastIndex = lastIndex + 1
			categoryIndex[category] = lastIndex
		}

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// Process chunks in parallel
		for g := 0; g < numGoroutines; g++ {
			start := g * chunkSize
			end := start + chunkSize
			if end > numElements {
				end = numElements
			}

			go func(start, end int) {
				defer wg.Done()
				for i := start; i < end; i++ {
					category := series.String[i]           // Get the category for this row
					columnIndex := categoryIndex[category] // Find the column index for the category
					df.Columns[columnIndex].Float[i] = 1.0 // Set the value to 1.0
				}
			}(start, end)
		}

		// Wait for all Goroutines to complete
		wg.Wait()
	}
	return nil
}

func (df *DataFrame) LabelEncode(identifiers ...any) error {
	var internalIndex int
	numGoroutines := runtime.NumCPU()
	length := df.GetLength()
	if length == 0 {
		return fmt.Errorf("dataframe is empty")
	}
	chunkSize := (length + numGoroutines - 1) / numGoroutines

	for _, identifier := range identifiers {
		series, err := df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("error retrieving column '%v': %w", identifier, err)
		}

		var equivalentMap map[interface{}]float64
		var nanLabel float64 = -1 // Special label for NaN values

		if series.DataType == "string" {
			uniqueValues := arrayUniqueValuesString(series.String)
			uniqueValues = parallelSortString(uniqueValues) // Sort in-place
			equivalentMap = make(map[interface{}]float64, len(uniqueValues))
			for index, value := range uniqueValues {
				equivalentMap[value] = float64(index)
			}
		} else {
			uniqueValues := arrayUniqueValuesFloat(series.Float)
			uniqueValues = ParallelSortFloat(uniqueValues) // Sort in-place
			equivalentMap = make(map[interface{}]float64, len(uniqueValues))
			internalIndex = 0
			for _, value := range uniqueValues {
				if math.IsNaN(value) {
					continue // Handle NaN separately
				}
				equivalentMap[value] = float64(internalIndex)
				internalIndex++
			}
		}

		// Concurrently process chunks to apply label encoding
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for g := 0; g < numGoroutines; g++ {
			start := g * chunkSize
			end := start + chunkSize
			if end > length {
				end = length
			}

			go func(start, end int) {
				defer wg.Done()
				for i := start; i < end; i++ {
					if series.DataType == "string" {
						series.Float[i] = equivalentMap[series.String[i]]
					} else {
						value := series.Float[i]
						if math.IsNaN(value) { // Handle NaN explicitly
							series.Float[i] = nanLabel
						} else {
							series.Float[i] = equivalentMap[value]
						}
					}
				}
			}(start, end)
		}

		wg.Wait()

		// Clear the original string data and update metadata
		if series.DataType == "string" {
			series.String = nil
		}
		series.DataType = "float"
	}

	return nil
}

/*
###################
#Feature Selection#
###################
*/

func PearsonCorrelation(x, y []float64) (float64, error) {
	xMean := arrayMean(x)
	yMean := arrayMean(y)

	numeratorChan := make(chan float64)
	xDenominatorChan := make(chan float64)
	yDenominatorChan := make(chan float64)

	// Compute numerator and denominators concurrently
	go func() {
		sum := 0.0
		for i := 0; i < len(x); i++ {
			xDiff := x[i] - xMean
			yDiff := y[i] - yMean
			sum += xDiff * yDiff
		}
		numeratorChan <- sum
	}()
	go func() {
		sum := 0.0
		for i := 0; i < len(x); i++ {
			xDiff := x[i] - xMean
			sum += math.Pow(xDiff, 2)
		}
		xDenominatorChan <- sum
	}()
	go func() {
		sum := 0.0
		for i := 0; i < len(y); i++ {
			yDiff := y[i] - yMean
			sum += math.Pow(yDiff, 2)
		}
		yDenominatorChan <- sum
	}()

	// Collect results from channels
	numerator := <-numeratorChan
	xDenominator := <-xDenominatorChan
	yDenominator := <-yDenominatorChan

	// Compute denominator
	denominator := math.Sqrt(xDenominator * yDenominator)
	if denominator == 0 {
		return 0, errors.New("division by zero in correlation calculation")
	}

	// Compute and return Pearson correlation
	return numerator / denominator, nil

}

func (df *DataFrame) SelectByCorrelation(targetIdentifier any, threshold float64) error {
	var targetSeries *Series
	var err error
	var selectedColumns []string
	var correlation float64

	targetSeries, err = df.GetColumnDynamic(targetIdentifier)
	if err != nil {
		return fmt.Errorf("error retrieving column '%v': %w", targetIdentifier, err)
	}
	if targetSeries.DataType != "float" {
		return fmt.Errorf("column '%v' is not a float column", targetIdentifier)
	}

	for _, column := range df.Columns {
		if column.Name == targetSeries.Name {
			continue
		}
		if column.DataType != "float" {
			continue
		}
		correlation, err = PearsonCorrelation(column.Float, targetSeries.Float)
		if err != nil {
			return fmt.Errorf("error calculating correlation for column '%v': %w", targetIdentifier, err)
		}

		if math.Abs(correlation) >= threshold {
			selectedColumns = append(selectedColumns, column.Name)
		}
	}

	var filteredColumns []Series
	for _, column := range df.Columns {
		if arrayContainsString(selectedColumns, column.Name) || column.DataType != "string" {
			filteredColumns = append(filteredColumns, column)
		}
	}
	df.Columns = filteredColumns
	return nil
}

func (df *DataFrame) VarianceThreshold(threshold float64) error {
	if df.GetLength() == 0 {
		return fmt.Errorf("dataframe is empty")
	}
	var selectedColumns []string
	var variance float64

	for _, column := range df.Columns {
		if column.DataType != "float" {
			continue
		}

		variance = arrayVariance(column.Float)
		if variance >= threshold {
			selectedColumns = append(selectedColumns, column.Name)
		}
	}

	var filteredColumns []Series
	for _, column := range df.Columns {
		if arrayContainsString(selectedColumns, column.Name) || column.DataType != "string" {
			filteredColumns = append(filteredColumns, column)
		}
	}
	df.Columns = filteredColumns
	return nil
}

/*
###########
#Utilities#
###########
*/
func TrainTestSplit(df DataFrame, testSize float64, randomState int) (DataFrame, DataFrame, error) {
	// Validate inputs
	if testSize < 0 || testSize > 1 {
		return DataFrame{}, DataFrame{}, fmt.Errorf("testSize must be between 0 and 1")
	}
	var err error
	numRows := df.GetLength()

	// Calculate the split index
	numTest := int(float64(numRows) * testSize)
	numTrain := numRows - numTest

	// Create a new random source
	var rng *rand.Rand
	if randomState != 0 {
		rng = rand.New(rand.NewSource(int64(randomState)))
	} else {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Generate a shuffled index
	indices := rng.Perm(numRows)

	// Create train and test indices
	trainIndices := indices[:numTrain]
	testIndices := indices[numTrain:]

	// Create train and test DataFrames
	trainSet, err := df.SelectRows(trainIndices)
	if err != nil {
		return DataFrame{}, DataFrame{}, err
	}
	testSet, err := df.SelectRows(testIndices)
	if err != nil {
		return DataFrame{}, DataFrame{}, err
	}

	return trainSet, testSet, nil

}
