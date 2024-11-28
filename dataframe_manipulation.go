package grizzly

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func (df *DataFrame) FilterFloat(identifier any, condition func(value float64) bool) error {
	var series *Series
	var err error
	series, err = df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to retrieve column to filter float %v: %w", identifier, err)
	}

	if series.DataType != "float" {
		return fmt.Errorf("column %v is not of type float; actual type is %q", identifier, series.DataType)
	}

	length := len(series.Float) // Use the actual length of the float slice
	if length == 0 {
		return nil
	}

	// Determine number of goroutines
	numGoroutines := runtime.NumCPU()
	chunkSize := (length + numGoroutines - 1) / numGoroutines // Calculate chunk size
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= length {
			break // Ensure we don't start beyond the slice length
		}
		if end > length {
			end = length // Adjust end index to stay within bounds
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if j >= length {
					// Double-check bounds to prevent any unexpected issues
					break
				}
				if condition(series.Float[j]) {
					for i, column := range df.Columns {
						if column.DataType == "float" {
							df.Columns[i].Float = append(column.Float[:j], column.Float[j+1:]...)
						} else {
							df.Columns[i].String = append(column.String[:j], column.String[j+1:]...)
						}
					}
				}
			}
		}(start, end)
	}

	// Closing channel after all goroutines finish
	go func() {
		wg.Wait()
	}()
	return nil
}

func (df *DataFrame) FilterString(identifier any, condition func(value string) bool) error {
	var series *Series
	var err error
	series, err = df.GetColumnDynamic(identifier)

	if err != nil {
		return fmt.Errorf("failed to retrieve column to filter string %v: %w", identifier, err)
	}

	if series.DataType != "float" {
		return fmt.Errorf("column %v is not of type string; actual type is %q", identifier, series.DataType)
	}

	length := len(series.Float) // Use the actual length of the float slice
	if length == 0 {
		return nil
	}

	// Determine number of goroutines
	numGoroutines := runtime.NumCPU()
	chunkSize := (length + numGoroutines - 1) / numGoroutines // Calculate chunk size
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= length {
			break // Ensure we don't start beyond the slice length
		}
		if end > length {
			end = length // Adjust end index to stay within bounds
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if j >= length {
					// Double-check bounds to prevent any unexpected issues
					break
				}
				if condition(series.String[j]) {
					for i, column := range df.Columns {
						if column.DataType == "float" {
							df.Columns[i].Float = append(column.Float[:j], column.Float[j+1:]...)
						} else {
							df.Columns[i].String = append(column.String[:j], column.String[j+1:]...)
						}
					}
				}
			}
		}(start, end)
	}

	// Closing channel after all goroutines finish
	go func() {
		wg.Wait()
	}()
	return nil
}

func (df *DataFrame) ApplyFloat(identifier any, operation func(float64) float64) error {
	var err error
	// Retrieve the series
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to retrieve column to apply float %v: %w", identifier, err)
	}

	if series.DataType != "float" {
		return fmt.Errorf("column %v is not of type float; actual type is %q", identifier, series.DataType)
	}

	// Get the length of the data
	numElements := len(series.Float)
	if numElements == 0 {
		return nil
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

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
				series.Float[i] = operation(series.Float[i])
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return nil
}

func (df *DataFrame) ApplyString(identifier any, operation func(string) string) error {
	var err error
	// Retrieve the series
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to retrieve column to apply string %v: %w", identifier, err)
	}

	if series.DataType != "string" {
		return fmt.Errorf("column %v is not of type string; actual type is %q", identifier, series.DataType)
	}

	// Get the length of the data
	numElements := len(series.String)
	if numElements == 0 {
		return nil
	}

	// Determine the number of goroutines based on available CPUs
	numGoroutines := runtime.NumCPU()
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

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
				series.String[i] = operation(series.String[i])
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return nil
}

func (df *DataFrame) ReplaceWholeWord(identifier any, old, new string) error {
	// Retrieve the column by name
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to replace whole word in column %v: %w", identifier, err)
	}

	// Perform the replacement
	series.ReplaceWholeWord(old, new)
	return nil
}

func (df *DataFrame) Replace(identifier any, old, new any) error {
	// Retrieve the column by name
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to replace in column %v: %w", identifier, err)
	}
	oldString, err := InterfaceConvertToString(old)
	newString, err := InterfaceConvertToString(new)
	if err != nil {
		return fmt.Errorf("failed to replace in column %v: %w", identifier, err)
	}

	// Perform the replacement
	series.Replace(oldString, newString)
	return nil
}

func (df *DataFrame) DropByIndex(index ...int) {
	var newSeries []Series
	for i, series := range df.Columns {
		if !ArrayContainsInteger(index, i) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
	return
}

func (df *DataFrame) DropByName(name ...string) {
	var newSeries []Series
	for _, series := range df.Columns {
		if !ArrayContainsString(name, series.Name) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
	return
}

func (df *DataFrame) DropDynamic(identifier any) error {
	var possibleName string
	var possibleIndex int
	var byIndex bool
	var err error

	switch v := identifier.(type) {
	case int:
		byIndex = true
		possibleIndex = v
	default:
		byIndex = false
		possibleName, err = InterfaceConvertToString(identifier)
		if err != nil {
			return err
		}
	}
	if byIndex {
		df.DropByIndex(possibleIndex)
		return nil
	}
	df.DropByName(possibleName)
	return nil
}

func (df *DataFrame) ConvertStringToFloat(identifiers ...any) error {
	var check *Series
	var err error

	for _, identifier := range identifiers {
		check, err = df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("failed to convert column %v from string to float: %w", identifier, err)
		}
		check.ConvertStringToFloat()
	}
	return nil
}

func (df *DataFrame) ConvertFloatToString(identifiers ...any) error {
	var check *Series
	var err error

	for _, identifier := range identifiers {
		check, err = df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("failed to convert column %v from string to float: %w", identifier, err)
		}
		check.ConvertFloatToString()
	}
	return nil
}

func (df *DataFrame) SplitColumn(identifier any, delimiter string, newColumnNames []string) error {
	var err error
	var column *Series
	column, err = df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to split column %v: %w", identifier, err)
	}

	if column.DataType == "float" {
		return fmt.Errorf("column %v cannot be split, select string column to split", identifier)
	}
	if len(newColumnNames) == 0 {
		return nil
	}

	numElements := column.GetLength()
	numGoroutines := runtime.NumCPU() // Use number of available CPUs for parallelism
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	// Create slices to hold the new column values
	splitValues := make([][]string, len(newColumnNames))
	for i := range splitValues {
		splitValues[i] = make([]string, numElements)
	}

	var wg sync.WaitGroup

	// Split the work among goroutines
	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numElements {
			end = numElements
		}
		if start >= end {
			break // No more elements to process
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				parts := strings.Split(column.String[i], delimiter)
				for j := 0; j < len(newColumnNames); j++ {
					if j < len(parts) {
						splitValues[j][i] = parts[j]
					} else {
						splitValues[j][i] = "" // Fill missing values with an empty string
					}
				}
			}
		}(start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Add new columns to the DataFrame
	for j, newName := range newColumnNames {
		newColumn := NewStringSeries(newName, splitValues[j])
		err = df.AddSeries(newColumn)
		if err != nil {
			return fmt.Errorf("failed to split column %q: %w", newName, err)
		}
	}

	return nil
}

func (df *DataFrame) JoinColumns(identifier1, identifier2 any, delimiter, newColumnName string) error {
	var column1 *Series
	var column2 *Series
	var err error
	// Retrieve the columns to be joined
	column1, err = df.GetColumnDynamic(identifier1)
	if err != nil {
		return fmt.Errorf("failed to join columns %v: %w", identifier1, err)
	}
	column2, err = df.GetColumnDynamic(identifier2)
	if err != nil {
		return fmt.Errorf("failed to join columns %v: %w", identifier2, err)
	}

	// Validate that both columns are string columns
	if column1.DataType != "string" || column2.DataType != "string" {
		return fmt.Errorf("JoinColumns is only supported for string columns")
	}

	numElements := column1.GetLength()
	joinedValues := make([]string, numElements)

	// Use goroutines to join columns in parallel for large datasets
	numGoroutines := runtime.NumCPU()
	if numGoroutines > numElements {
		numGoroutines = numElements // Limit the number of goroutines to the number of elements
	}
	chunkSize := (numElements + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if end > numElements {
			end = numElements
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				joinedValues[i] = column1.String[i] + delimiter + column2.String[i]
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Add the new joined column to the DataFrame
	newColumn := NewStringSeries(newColumnName, joinedValues)
	err = df.AddSeries(newColumn)
	if err != nil {
		return fmt.Errorf("failed to join columns %q: %w", newColumnName, err)
	}
	return nil
}

func (df *DataFrame) SliceRows(low int, high int) error {
	if low < 0 || high >= df.GetNumberOfColumns() {
		return fmt.Errorf("out of range")
	}
	for i := range df.Columns {
		if df.Columns[i].DataType == "float" {
			df.Columns[i].Float = df.Columns[i].Float[low:high]
		} else {
			df.Columns[i].String = df.Columns[i].String[low:high]
		}
	}
	return nil
}

func (df *DataFrame) SelectRows(indices []int) (DataFrame, error) {
	// Check if the indices are valid
	for _, index := range indices {
		if index < 0 || index >= df.GetLength() {
			return DataFrame{}, fmt.Errorf("index out of range: %d", index)
		}
	}

	// Create a new DataFrame to store the selected rows
	selected := DataFrame{
		Columns: make([]Series, len(df.Columns)),
	}

	// Iterate over each column to copy selected rows
	for i, col := range df.Columns {
		selected.Columns[i] = Series{
			Name:     col.Name,
			DataType: col.DataType,
		}
		if col.DataType == "float" {
			for _, index := range indices {
				selected.Columns[i].Float = append(selected.Columns[i].Float, col.Float[index])
			}
		} else { // For "string" or other types
			for _, index := range indices {
				selected.Columns[i].String = append(selected.Columns[i].String, col.String[index])
			}
		}
	}

	return selected, nil
}

func (df *DataFrame) SliceColumns(low, high int) error {
	if low < 0 || high >= df.GetNumberOfColumns() {
		return fmt.Errorf("out of range")
	}
	df.Columns = df.Columns[low:high]
	return nil
}

func (df *DataFrame) MergeDataFrame(otherDf DataFrame) error {
	names := df.GetColumnNames()
	otherNames := otherDf.GetColumnNames()
	for _, name := range names {
		if ArrayContainsString(otherNames, name) {
			return fmt.Errorf("column already exists")
		}
	}
	for _, column := range otherDf.Columns {
		df.AddSeriesForced(column)
	}
	return nil
}

func (df *DataFrame) Concatenate(otherDf DataFrame) error {
	var newColumn Series
	var err error

	otherNames := otherDf.GetColumnNames()
	names := df.GetColumnNames()
	for _, name := range otherNames {
		if !ArrayContainsString(names, name) {
			newColumn = NewStringSeries(name, []string{})
			df.AddSeriesForced(newColumn)
		}
	}
	names = append(names, otherNames...)
	var series *Series
	var otherSeries *Series
	for _, name := range names {
		series, err = df.GetColumnByName(name)
		if err != nil {
			return fmt.Errorf("failed to concatenate dataframes %q: %w", name, err)
		}
		otherSeries, err = otherDf.GetColumnByName(name)
		if err != nil {
			return fmt.Errorf("failed to concatenate dataframes %q: %w", name, err)
		}
		if series.DataType == "float" && otherSeries.DataType == "string" {
			err = df.ConvertFloatToString(name)
			if err != nil {
				return fmt.Errorf("failed to concatenate dataframes %q: %w", name, err)
			}
		} else if series.DataType == "string" && otherSeries.DataType == "float" {
			err = df.ConvertStringToFloat(name)
			if err != nil {
				return fmt.Errorf("failed to concatenate dataframes %q: %w", name, err)
			}
		}

		if series.DataType == "float" {
			series.Float = append(series.Float, otherSeries.Float...)
		} else if series.DataType == "string" {
			series.String = append(series.String, otherSeries.String...)
		}
	}
	return nil
}

func (df *DataFrame) DuplicateColumn(identifiers ...any) error {
	var ptr *Series
	var series Series
	var err error

	for _, identifier := range identifiers {
		ptr, err = df.GetColumnDynamic(identifier)
		if err != nil {
			return fmt.Errorf("failed to duplicate column %v: %w", identifier, err)
		}
		series = *ptr
		series.Name = series.Name + strconv.Itoa(1)
		err = df.AddSeries(series)
		if err != nil {
			return fmt.Errorf("failed to duplicate column %v: %w", identifier, err)
		}
	}
	return nil
}

func (df *DataFrame) MathBase(identifier1, identifier2 any, newColumnName string, operation func(float64, float64) float64) error {
	var newColumn Series
	var series1 *Series
	var series2 *Series
	var err error

	series1, err = df.GetColumnDynamic(identifier1)
	if err != nil {
		return fmt.Errorf("failed to execute math operation %v: %w", identifier1, err)
	}
	series2, err = df.GetColumnDynamic(identifier2)
	if err != nil {
		return fmt.Errorf("failed to execute math operation %v: %w", identifier2, err)
	}
	if !(series1.DataType == "float") || !(series2.DataType == "float") {
		return fmt.Errorf("math operation only supports floating point values")
	}
	size := series1.GetLength()
	if size == 0 {
		return nil
	}
	newColumn = Series{
		Name:     newColumnName,
		DataType: "float",           // Default type
		String:   make([]string, 0), // Preallocate with length
		Float:    make([]float64, size),
	}

	numGoroutines := runtime.NumCPU()
	chunkSize := (size + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup

	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if start >= size {
			break // Ensure we don't start beyond the slice length
		}
		if end > size {
			end = size // Adjust end index for the last chunk
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				newColumn.Float[j] = operation(series1.Float[j], series2.Float[j])
			}
		}(start, end)
	}

	wg.Wait()

	err = df.AddSeries(newColumn)
	if err != nil {
		return fmt.Errorf("failed to execute math operation %q: %w", newColumnName, err)
	}
	return nil
}

func (df *DataFrame) Sum(identifier1, identifier2 any, newColumnName string) error {
	return df.MathBase(identifier1, identifier2, newColumnName, func(x, y float64) float64 { return x + y })
}

func (df *DataFrame) Subtraction(identifier1, identifier2 any, newColumnName string) error {
	return df.MathBase(identifier1, identifier2, newColumnName, func(x, y float64) float64 { return x - y })
}

func (df *DataFrame) Multiplication(identifier1, identifier2 any, newColumnName string) error {
	return df.MathBase(identifier1, identifier2, newColumnName, func(x, y float64) float64 { return x * y })
}

func (df *DataFrame) Division(identifier1, identifier2 any, newColumnName string) error {
	return df.MathBase(identifier1, identifier2, newColumnName, func(x, y float64) float64 { return x / y })
}

func (df *DataFrame) SetFloatValue(identifier any, rowIndex int, newValue float64) error {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to set value for column %v: %w", identifier, err)
	}
	if series.DataType != "float" {
		return fmt.Errorf("column %v only supports floating point values", identifier)
	}
	series.Float[rowIndex] = newValue
	return nil
}

func (df *DataFrame) SetStringValue(identifier any, rowIndex int, newValue string) error {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to set value for column %v: %w", identifier, err)
	}
	if series.DataType != "string" {
		return fmt.Errorf("column %v only supports string values", identifier)
	}
	series.String[rowIndex] = newValue
	return nil
}

func (df *DataFrame) SetValue(identifier any, rowIndex int, newValue any) error {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("failed to set value for column %v: %w", identifier, err)
	}
	if series.DataType == "string" {
		newS, err := InterfaceConvertToString(newValue)
		if err != nil {
			return fmt.Errorf("failed to set value for column %v: %w", identifier, err)
		}
		series.String[rowIndex] = newS
		return nil
	}
	newF, err := InterfaceConvertToFloat(newValue)
	if err != nil {
		return fmt.Errorf("failed to set value for column %v: %w", identifier, err)
	}
	series.Float[rowIndex] = newF
	return nil
}

func (df *DataFrame) GetFloatValue(identifier any, rowIndex int) (float64, error) {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return 0, fmt.Errorf("failed to get value for column %v: %w", identifier, err)
	}
	if series.DataType != "float" {
		return 0, fmt.Errorf("column %q only supports floating point tasks", identifier)
	}
	return series.Float[rowIndex], nil
}

func (df *DataFrame) GetStringValue(identifier, rowIndex int) (string, error) {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to get value for column %v: %w", identifier, err)
	}
	if series.DataType != "string" {
		return "", fmt.Errorf("column %q only supports string tasks", identifier)
	}
	return series.String[rowIndex], nil
}

func (df *DataFrame) GetValue(identifier any, rowIndex int) (any, error) {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get value for column %v: %w", identifier, err)
	}
	if series.DataType == "float" {
		return series.Float[rowIndex], nil
	}
	return series.String[rowIndex], nil
}

func (df *DataFrame) Expand(size int, defaultFloat float64, defaultString string) {
	for i, series := range df.Columns {
		if series.DataType == "string" {
			temp := make([]string, size)
			for n := range temp {
				temp[n] = defaultString
			}
			df.Columns[i].String = append(series.String, temp...)
			return
		}
		temp := make([]float64, size)
		for n := range temp {
			temp[n] = defaultFloat
		}
		df.Columns[i].Float = append(series.Float, temp...)
		return
	}
}

func (df *DataFrame) SwapRows(index1, index2 int) error {
	size := df.GetLength()
	if index1 < 0 || index1 >= size || index2 < 0 || index2 >= size {
		return fmt.Errorf("row index out of bounds")
	}
	if index1 == index2 {
		return nil
	}
	for i, series := range df.Columns {
		if series.DataType == "float" {
			df.Columns[i].Float[index1], df.Columns[i].Float[index2] = df.Columns[i].Float[index2], df.Columns[i].Float[index1]
		} else {
			df.Columns[i].String[index1], df.Columns[i].String[index2] = df.Columns[i].String[index2], df.Columns[i].String[index1]
		}
	}
	return nil
}

// Sort QuickSort sorts the array in place using the QuickSort algorithm
func (df *DataFrame) Sort(identifier any, internal ...int) error {
	var low int
	var high int
	var p int
	var err error
	if len(internal) == 0 {
		low = 0
		high = df.GetLength() - 1
	} else {
		low = internal[0]
		high = internal[1]
	}
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return fmt.Errorf("error sorting %v: %w", identifier, err)
	}
	if low < high {
		// Partition the array
		if series.DataType == "float" {
			p, err = df.PartitionFloat(series, low, high)
			if err != nil {
				return fmt.Errorf("error sorting dataframe by %q column", series.Name)
			}
		} else {
			p, err = df.PartitionString(series, low, high)
			if err != nil {
				return fmt.Errorf("error sorting dataframe by %q column", series.Name)
			}
		}

		// Recursively sort the sub-arrays
		err = df.Sort(identifier, low, p-1)
		if err != nil {
			return fmt.Errorf("error sorting dataframe")
		}
		err = df.Sort(identifier, p+1, high)
		if err != nil {
			return fmt.Errorf("error sorting dataframe")
		}
	}
	return nil
}

// PartitionString rearranges the array and returns the pivot index for strings
func (df *DataFrame) PartitionString(column *Series, low, high int) (int, error) {
	var err error
	pivot := column.String[high] // Choose the last element as pivot
	i := low - 1                 // Pointer for the smaller element

	for j := low; j < high; j++ {
		if column.String[j] < pivot {
			i++
			err = df.SwapRows(i, j) // Swap if element is smaller than pivot
			if err != nil {
				return i, fmt.Errorf("failed to sort df: %w", err)
			}
		}
	}

	// Place the pivot in the correct position
	err = df.SwapRows(i+1, high)
	if err != nil {
		return i, fmt.Errorf("failed to sort df: %w", err)
	}
	return i + 1, nil
}

// PartitionFloat partition rearranges the array and returns the pivot index
func (df *DataFrame) PartitionFloat(column *Series, low, high int) (int, error) {
	var err error

	pivot := column.Float[high] // Choose the last element as pivot
	i := low - 1                // Pointer for the smaller element

	for j := low; j < high; j++ {
		if column.Float[j] < pivot {
			i++
			err = df.SwapRows(i, j) // Swap if element is smaller than pivot
			if err != nil {
				return i, fmt.Errorf("failed to sort df: %w", err)
			}
		}
	}

	// Place the pivot in the correct position
	err = df.SwapRows(i+1, high)
	if err != nil {
		return i, fmt.Errorf("failed to sort df: %w", err)
	}
	return i + 1, nil
}
