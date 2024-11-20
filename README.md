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
