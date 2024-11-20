# Go DataFrame Library
An alternative to Python's Pandas library, this package provides efficient data manipulation and aggregation capabilities tailored for Go developers. The library focuses on DataFrame operations with support for series operations as backend functionalities.

### Features
- Flexible DataFrame creation and manipulation
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
- otherDf *DataFrame*: the new columns will extract from it.
- defaultValue *string*: default value for empty values.
