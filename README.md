# Go DataFrame Library
Grizzly is a DataFrame library for Go, designed to harness the full power of GoRoutines for handling large datasets efficiently. Its core aim is to provide an easy-to-use, yet robust, solution for data manipulation while maximizing the computational capabilities of modern machines through parallelized task execution.

Unlike many other libraries, Grizzly enforces a more rigid approach to DataFrame management. Users are required to explicitly specify data types, such as float or string, ensuring clarity, type safety, and reducing potential errors in data processing.

![image](https://github.com/user-attachments/assets/8e8ed677-ee0c-4c13-9cf0-b6c48b009da6)


### Features
- Typed DataFrame (float64 and string)
- Data aggregation and statistical functions
- Data attributes and metadata handling
- Import and export utilities
  
### Installation
To install the package, use:
```
go get github.com/Puchungualotsqui/grizzly
```
If you have problems. Please, try to download directly the package.
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
max, _ = df.GetMax()
```
### GetMin
Return a DataFrame with the min of each column.
```
var min DataFrame
min, _ = df.GetMin()
```
### GetMean
Return a DataFrame with the mean of each column.
```
var mean DataFrame
mean, _ = df.GetMean()
```
### GetMedian
Return a DataFrame with the median of each column.
```
var median DataFrame
median, _ = df.GetMedian()
```
### GetProduct
Return a DataFrame with the product of each column.
```
var product DataFrame
product, _ = df.GetProduct()
```
### GetSum
Return a DataFrame with the summation of each column.
```
var sum DataFrame
sum, _ = df.GetSum()
```
### GetVariance
Return a DataFrame with the variance of each column.
```
var variance DataFrame
variance, _ = df.GetVariance()
```
### CountWord
Return a DataFrame with the count of the input string.
- word *string*: word to count.
```
var count DataFrame
count, _ = df.CountWord("hello")
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
var uniqueValues DataFrame
uniqueValues = df.GetUniqueValues()
```
### CountNaNValues
Return a DataFrame with the NaN values of each column.
```
var nanCount DataFrame
nanCount = df.CountNaNValues()
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
### GetColumnByName
Return a reference to the column with the given name.
- name *string*: name of the column to return.
```
var column *Series
column, _ = df.GetColumnByName("names")
```
### GetColumnByIndex
Return a reference to the column of the given index.
- index *int*: index of the asked column.
```
var column *Series
column, _ = df.GetColumnByIndex(4)
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
dataType, _ = df.GetColumnType("name")
```
### ColumnIsStringIndex
Return a bool true or false depending if the data type of the selected column is string.
- index *int*: the index of the column to verify if is string type.
```
var isString bool
isString, _ = df.ColumnIsStringIndex(0)
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
isFloat, _ = df.ColumnIsFloat(0)
```
### GetColumnIndexByName
Return the index of the column by name.
- columnName *string*: column name
```
var index int
index, _ = df.GetColumnIndexByName("name")
```
## DataFrame Manipulation
### FilterFloat
Filter rows based on a condition for float columns.
- identifier *any*: integer or name of the column to filter.
- condition *func(value float64) bool*: function to filter. True values will be deleted.
```
filter := func(val float64) bool {
	if val >= 4 {
		return true
	}
	return false
}
df.FilterFloat(0, filter)
```
### FilterString
Filter rows based on a condition for string columns.
- identifier *any*: integer or name of the column to filter.
- condition *func(value string) bool*: function to filter. True values will be deleted.
```
filter := func(val string) bool {
	if val == "david" {
		return true
	}
	return false
}
df.FilterFloat(0, filter)
```
### ApplyFloat
Apply a function to transform a float column.
- identifier *any*: integer or name of the column to apply operation.
- operation *func(float64) float64*: operation to apply over the column.
```
isWeekend := func(val float64) float64 {
	return val * 2
	}
df.ApplyString("is_weekend", isWeekend)
```
### ApplyString
Apply a function to transform a string column.
- identifier *any*: integer or name of the column to apply operation.
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
- identifier *any*: integer or name of the column to apply change.
- old *string*: word to search to replace.
- new *string*: replacement of the old string.
```
df.ReplaceWholeWord("name", "Dabid", "David")
```
### Replace
Replace any substring equal to old value.
Replace the whole value of each value.
- identifier *any*: integer or name of the column to apply change.
- old *any*: substring or number to search to replace.
- new *any*: replacement of the old value.
```
df.Replace("name", "Dabid", "David")
```
### DropByIndex
Drop column indicating the index.
- index *...int*: index to drop.
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
- identifiers *...any*: name or index of the column to convert to float.
```
df.ConvertStringToFloat("ages", "salary")
```
### ConvertFloatToString
Try to convert a float column into string column.
- identifiers *...any*: name or index of the column to convert to string.
```
df.ConvertFloatToString("postal_code")
```
### SplitColumn
Split string column in different columns based on a delimiter.
- identifier *any*: name or index of the column to split the data.
- delimiter *string*: delimiter to split the data.
- newColumnNames *[]string*: names for the new columns.
```
df.SplitColumn("date_sold", "/", []string{"day", "month", "year"})
```
### JoinColumns
Combine two string columns in a DataFrame by joining their values with a specified delimiter and creates a new column with the resulting values.
- identifier1 *any*: name or index of the left column to join data.
- identifier2 *any*: name or index the right column to join data.
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
### SliceColumns
Slice the columns based on index number.
- low *int*: initial index to slice.
- high *int*: final index to slice.
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
- identifier *...nay*: names or indexes of the columns to duplicate.
```
df.DuplicateColumn("names","street")
```
### Sum
Sum two columns, the result is saved in a new column.
- identifier1 *any*: first column to sum.
- identifier2 *any*: second column to sum.
- newColumnName *string*: name for the new column with the result.
```
df.Sum("x","y","x+y")
```
### Subtraction
Subtract two columns, the result is saved in a new column.
- identifier1 *any*: column to be minuend.
- identifier2 *any*: column to be subtrahen.
- newColumnName *string*: name for the new column with the result.
```
df.Subtraction("x","y","x-y")
```
### Multiplication
Multiply two columns, the result is saved in a new column.
- identifier1 *any*: first column to multiply.
- identifier2 *any*: second column to multiply.
- newColumnName *string*: name for the new column with the result.
```
df.Subtraction("x","y","x*y")
```
### Division
Divide two columns, the result is saved in a new column.
- identifier1 *any*: column to be dividend.
- identifier2 *any*: column to be divisor.
- newColumnName *string*: name for the new column with the result.
```
df.Subtraction("x","y","x/y")
```
### SetFloatValue
Change the value of a row of a float column.
- identifier *any*: index or name of the column to change value.
- rowIndex *int*: index of the row to change the value.
- newValue *float64*: new value.
```
df.SetFloatValue(2,100,94.213)
```
### SetStringValue
Change the value of a row of a string column.
- identifier *any*: index or name of the column to change value.
- rowIndex *int*: index of the row to change the value.
- newValue *string*: new value.
```
df.SetStringValue(2,100,"sebastian")
```
### SetValue
Change the value of a row of an column. It will try to convert it to the type of the column.
- identifier *any*: index or name of the column to change value.
- rowIndex *int*: index of the row to change the value.
- newValue *any*: new value.
```
df.SetValue(2,100,143)
```
### GetFloatValue
Return the float64 value of a row of a float column.
- identifier *any*: index or name of the column to return value.
- rowIndex *int*: index of the row to return the value.
```
var returnValue float
returnValue = df.GetFloatValue(2,1245)
```
### GetStringValue
Return the string value of a row of a float column.
- identifier *any*: index or name of the column to return value.
- rowIndex *int*: index of the row to return the value.
```
var returnValue string
returnValue = df.GetStringValue(2,1245)
```
### GetValue
Return the value of a row of a column.
- identifier *any*: index or name of the column to return value.
- rowIndex *int*: index of the row to return the value.
```
var returnValue any
returnValue = df.GetValue(3,123123)
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
- identifier *any*: index or name of the column to sort the dataframe.
```
df.Sort("name")
```
## Data Cleaning
### FillNaN
Replace all NaN values in float columns.
- newValue *float64*: value to replace all NaN values.
- identifiers *...any*: name or indexes of columns to replace NaN values. If it is empty will replace NaN values in all columns.
```
df.FillNaN(5, "ages","salary")
```
### DropNaN
Drop all rows with NaN values in float columns.
- identifiers *...any*: name or indexes of columns to check if there is any NaN value. If it is empty will check all columns.
```
df.DropNaN()
```
### RemoveOutliersZScore
Remove rows with outliers values on a column. Using ZScore method.
- identifier *any*: name or index of the column to check outliers.
- threshold *float64*: break point to filter outliers.
```
df.RemoveOutliersZScore("salary", 0.25)
```
### RemoveOutliersIQR
Remove outliers using interquartile range.
- identifier *any*: name or index of the column to check outliers.
```
df.RemoveOutliersIQR("salary")
```
### RemoveDuplicates
Remove rows with exact same values.
```
df.RemoveDuplicates()
```
## Input
### ImportCSV
Import CSV file as Grizzly DataFrame.
- filepath *string*: file path of the csv file.
```
df, _ = grizzly.ImportCSV("example.csv")
```
## Output
### ExportToCSV
Export Dataframe as a CSV file.
- filepath *string*: file path for the csv file.
```
df.ExportToCSV("example.csv")
```
### ExportToCSVSimple
Export Dataframe as a CSV file. It is better for small Dataframes.
```
df.ExportToCSVSimple("example.csv")
```
