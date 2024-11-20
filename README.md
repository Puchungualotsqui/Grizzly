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
    filePath := "example.csv"
    df, err := grizzly.ImportCSV(filePath)
}
```
### DataFrame Manipulation
Example of basic manipulation functions:
```
df.ReplaceWholeWord("column_name", "USA", "United States"

filterFloat := func(value float64) bool {
			return value > 25
		}
df.FilterFloat("column_name", filterFloat)
```
### Aggregation
Aggregate functions include:
```
mean := df.GetMean()
empty := df.CountEmpty()
```

# Functions
## DataFrame Creation
### CreateDataFrame
It initialize the DataFrame structure.
```
df := CreateDataFrame
```
