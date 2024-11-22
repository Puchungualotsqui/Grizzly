# Go DataFrame Library
Grizzly is a DataFrame library for Go, designed to harness the full power of GoRoutines for handling large datasets efficiently. Its core aim is to provide an easy-to-use, yet robust, solution for data manipulation while maximizing the computational capabilities of modern machines through parallelized task execution.

Unlike many other libraries, Grizzly enforces a more rigid approach to DataFrame management. Users are required to explicitly specify data types, such as float or string, ensuring clarity, type safety, and reducing potential errors in data processing.

![image](https://github.com/user-attachments/assets/8e8ed677-ee0c-4c13-9cf0-b6c48b009da6)


### Features
- Typed DataFrame (float64 and string)
- Data aggregation and statistical functions
- Data attributes and metadata handling
- Import and export utilities for common data formats
  
### Installation
To install the package, use:
```
go get github.com/Puchungualotsqui/grizzly
```
## Basic Usage
### Creating a DataFrame

```
package main

import (
	"fmt"
	"grizzly"
)

func main() {
	names := []string{"Alice", "Bob", "Charlie", "Diana", "Ethan"}
	ages := []float64{25, 30, 35, 40, 28}

	// Initialize the DataFrame
	df := grizzly.CreateDataFrame()

	// Add columns
	df.CreateStringColumn("Names", names)
	df.CreateFloatColumn("Ages", ages)

	df.PrintHead(5)
}
```
# Functions
## DataFrame Creation
### CreateDataFrame
Initialize the DataFrame structure. All the modification should be done over the DataFrame structure.
```
df := CreateDataFrame()
```
### CreateFloatColumn
Create new float column.
- name *string*: name of column.
- nums *[]float64*: data of the column in float type.
```
ages := []float64{25, 30, 35, 40, 28}
df.CreateFloatColumn("Ages", ages)
```
### CreateStringColumn
Create new string column.
- name *string*: name of column.
- words *[]float64*: data of the column in string type.
```
names := []string{"Alice", "Bob", "Charlie", "Diana", "Ethan"}
df.CreateFloatColumn("Names", names)
```
### Print
Prints data in console.
- min *int*: starting index to print.
- max *int*: end index to print.
```
df.Print(5,10)
```
### PrintHead
Print first rows in console.
- max *int*: end index to print.
```
df.PrintHead(5)
```
### PrintTail
Print last rows in console.
- min *int*: start index to print.
```
df.PrintTail(5)
```
## DataFrame Aggregation
### GetMax
Return a DataFrame with the max of each column.
```
var max DataFrame
max = df.GetMax()
```
### GetMin
Return a DataFrame with the min of each column.
```
var min DataFrame
min = df.GetMin()
```
### GetMean
Return a DataFrame with the mean of each column.
```
var mean DataFrame
mean = df.GetMean()
```
### GetMedian
Return a DataFrame with the median of each column.
```
var median DataFrame
median = df.GetMedian()
```
### GetProduct
Return a DataFrame with the product of each column.
```
var product DataFrame
product = df.GetProduct()
```
### GetSum
Return a DataFrame with the summation of each column.
```
var sum DataFrame
sum = df.GetSum()
```
### GetVariance
Return a DataFrame with the variance of each column.
```
var variance DataFrame
variance = df.GetVariance()
```
### CountWord
Return a DataFrame with the count of the input string.
- word *string*: word to count.
```
var count DataFrame
count = df.CountWord("hello")
```
### GetNonFloatValues
Return a DataFrame with the non float values of each column.
```
var nonFloat DataFrame
nonFloat = df.GetNonFloatValues()
```
### GetUniqueValues
Return a DataFrame with the unique values of each column.
```
df.GetUniqueValues()
```
## DataFrame Attributes
### GetLength
Return the number of rows as integer.
```
var length int
length = df.GetLength()
```
### GetNumberOfColumns
Return the number of columns as integer.
```
var numColumns int
numColumns = df.GetNumberOfColumns()
```
### GetColumnNames
Return an array of strings with the names of the columns.
```
var columnNames []string
columnNames = df.GetColumnNames()
```
### GetShape
Return an array with the [number of columns, number of rows].
```
var shape []int
shape = df.GetShape()
```
### ContainsColumn
Return bool value as true if the column exists.
- name *string*: name of the column to verify.
```
var contains bool
contains = df.ContainsColumn("last_names")
```
### GetColumnTypeIndex
Return an string with the type of data of the column.
- index *int*: the index of the column to get the data type of it.
```
var dataType string
dataType = df.GetColumnTypeIndex(0)
```
### GetColumnType
Return an string with the type of data of the column.
- name *string*: the name of the column to get the data type of it.
```
var dataType string
dataType = df.GetColumnType("name")
```
### ColumnIsStringIndex
Return a bool true or false depending if the data type of the selected column is string.
- index *int*: the index of the column to verify if is string type.
```
var isString bool
isString = df.ColumnIsStringIndex(0)
```
### ColumnIsString
Return a bool true or false depending if the data type of the selected column is string.
- name *string*: the name of the column to verify if is string type.
```
var isString bool
isString = df.ColumnIsString(0)
```
### ColumnIsFloatIndex
Return a bool true or false depending if the data type of the selected column is float.
- index *int*: the index of the column to verify if is float type.
```
var isFloat bool
isFloat = df.ColumnIsFloatIndex(0)
```
### ColumnIsFloat
Return a bool true or false depending if the data type of the selected column is float.
- name *string*: the name of the column to verify if is float type.
```
var isFloat bool
isFloat = df.ColumnIsFloat(0)
```
## DataFrame Manipulation
### FilterFloat
Filter rows based on a condition for float columns.
- columnName *string*: name of the column for the filter.
- condition *func(value float64) bool*: function to filter. True values will be deleted.
```
filter := func(val float64) bool {
	if val >= 4 {
		return true
	}
	return false
}
df.FilterFloat("num_bath", filter)
```
### FilterString
Filter rows based on a condition for string columns.
- columnName *string*: name of the column for the filter.
- condition *func(value string) bool*: function to filter. True values will be deleted.
```
filter := func(val string) bool {
	if val == "david" {
		return true
	}
	return false
}
df.FilterFloat("names", filter)
```
### ApplyFloat
Apply a function to transform a float column.
- columnName *string*: name of the column to transform.
- operation *func(float64) float64*: operation to apply over the column.
```
isWeekend := func(val float64) float64 {
	return val * 2
	}
df.ApplyString("is_weekend", isWeekend)
```
### ApplyString
Apply a function to transform a string column.
- columnName *string*: name of the column to transform.
- operation *func(string) string*: operation to apply over the column.
```
isWeekend := func(val string) string {
	parsedDate, _ := time.Parse("2015-01-16", val)
	day := parsedDate.Weekday()
	if day == time.Saturday || day == time.Sunday {
		return "1"
	}
	return "0"
}
df.ApplyString("is_weekend", isWeekend)
```
### ReplaceWholeWord
Replace the whole value of each value.
- columnName *string*: name of the column to replace the data.
- old *string*: word to search to replace.
- new *string*: replacement of the old string.
```
df.ReplaceWholeWord("name", "Dabid", "David")
```
### Replace
Replace any substring equal to old value.
Replace the whole value of each value.
- columnName *string*: name of the column to replace the data.
- old *string*: substring to search to replace.
- new *string*: replacement of the old string.
```
df.Replace("name", "Dabid", "David")
```
### DropByIndex
Drop column indicating the index.
- index *..int*: index to drop.
```
df.DropByIndex(1,2,3,4)
```
### DropByName
Drop column indicating the names.
- name *...string*: column names to drop.
```
df.DropByName("names", "ages")
```
### ConvertStringToFloat
Try to convert a string column into float column.
- names *...string*: names of columns to convert columns.
```
df.ConvertStringToFloat("ages", "salary")
```
### ConvertFloatToString
Try to convert a float column into string column.
- names *...string*: names of columns to convert columns.
```
df.ConvertFloatToString("postal_code")
```
### ConvertStringToFloatIndex
Try to convert a string column into float column indicating the index of the columns.
- index *...int*: index of columns to convert columns.
```
df.ConvertStringToFloatIndex(2,5)
```
### ConvertFloatToStringIndex
Try to convert a float column into string column indicating the index of the columns.
- index *...int*: index of columns to convert columns.
```
df.ConvertFloatToStringIndex(2,5)
```
### SplitColumn
Split string column in different columns based on a delimiter.
- columnName *string*: name of the column to split the data.
- delimiter *string*: delimiter to split the data.
- newColumnNames *[]string*: names for the new columns.
```
df.SplitColumn("date_sold", "/", []string{"day", "month", "year"})
```
### JoinColumns
Combine two string columns in a DataFrame by joining their values with a specified delimiter and creates a new column with the resulting values.
- columnName1 *string*: name of the left column to join data.
- columnName2 *string*: name of the right column to join data.
- delimiter *string*: value between values.
- newColumnName *string*: name of the new column.
```
df.JoinColumns("date","time",":","full_date")
```
### SliceRows
Slice the rows of the dataframe.
- offset *int*: initial index to slice.
- length *int*: length of slice.
```
df.SliceRows(5,2)
```
### SliceColumnsByIndex
Slice the columns based on index number.
- offset *int*: initial index to slice.
- length *int*: length of slice.
```
df.SliceColumnsByIndex(5,2)
```
### MergeDataFrame
Add the columns of the other dataframe.
- otherDf *DataFrame*: the new columns will be extracted from it.
- defaultValue *string*: default value for empty values.
```
df.MergeDataFrame(otherDF, "None")
```
### Concatenate
Add the new columns of the other dataframe.
- otherDf *DataFrame*: the new rows will be extracted from it.
- defaultValue *string*: default value for empty values.
```
df.Concatenate(otherDF, "None")
```
### DuplicateColumn
Create a copy of a column.
- names *...string*: names of the columns to duplicate.
```
df.DuplicateColumn("names","street")
```
### Sum
Sum two columns, the result is saved in a new column.
- columnName1 *string*: first column to sum.
- columnName2 *string*: second column to sum.
- newColumnName *string*: name for the new column with the result.
```
df.Sum("x","y","x+y")
```
### Subtraction
Subtract two columns, the result is saved in a new column.
- columnName1 *string*: column to be minuend.
- columnName2 *string*: column to be subtrahen.
- newColumnName *string*: name for the new column with the result.
```
df.Subtraction("x","y","x-y")
```
### Multiplication
Multiply two columns, the result is saved in a new column.
- columnName1 *string*: first column to multiply.
- columnName2 *string*: second column to multiply.
- newColumnName *string*: name for the new column with the result.
```
df.Subtraction("x","y","x*y")
```
### Division
Divide two columns, the result is saved in a new column.
- columnName1 *string*: column to be dividend.
- columnName2 *string*: column to be divisor.
- newColumnName *string*: name for the new column with the result.
```
df.Subtraction("x","y","x/y")
```
### SetFloatValue
Change the value of a row of a float column.
- columnIndex *int*: index of the column to change value.
- rowIndex *int*: index of the row to change the value.
- newValue *float64*: new value.
```
df.SetFloatValue(2,100,94.213)
```
### SetStringValue
Change the value of a row of a string column.
- columnIndex *int*: index of the column to change value.
- rowIndex *int*: index of the row to change the value.
- newValue *string*: new value.
```
df.SetStringValue(2,100,"sebastian")
```
### GetFloatValue
Return the float64 value of a row of a float column.
- columnIndex *int*: index of the column to return value.
- rowIndex *int*: index of the row to return the value.
```
var returnValue float
returnValue = df.GetFloatValue(2,1245)
```
### GetStringValue
Return the string value of a row of a float column.
- columnIndex *int*: index of the column to return value.
- rowIndex *int*: index of the row to return the value.
```
var returnValue string
returnValue = df.GetStringValue(2,1245)
```
### Expand
Add new rows to the DataFrame.
- size *int*: amount of new rows.
- defaultFloat *float64*: default value for new rows in float columns.
- defaultString *string*: default value for new rows in string columns.
```
df.Expand(100, 0, "")
```
### SwapRows
Swap the value of two rows.
- index1 *int*: index of the first row to swap values.
- index2 *int*: index of the second row to swap values.
```
df.SwapRows(1,0)
```
### Sort
Sort the Dataframe based on one column.
- index *int*: index of the column to sort the dataframe.
```
df.Sort(0)
```
## Input
### ImportCSV
Import CSV file as Grizzly DataFrame.
- filepath *string*: file path of the csv file.
```
df = grizzly.ImportCSV("example.csv")
```
