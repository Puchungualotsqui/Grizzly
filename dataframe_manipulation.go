package grizzly

import "fmt"

func (df *DataFrame) GetColumnByName(name string) Series {
	for _, series := range df.Columns {
		if series.Name == name {
			return series
		}
	}
	panic(fmt.Sprintf("%s not found", name))
}

func (df *DataFrame) GetColumnByIndex(index int) Series {
	return df.Columns[index]
}

func (df *DataFrame) FilterFloat(seriesName string, condition func(float64) bool) {
	// Verify if series exists
	series := df.GetColumnByName(seriesName)
	if series.Name == "" {
		fmt.Println("Column not found")
	} else if series.DataType != "float" {
		fmt.Println("Not a float")
	} else {
		indexes := series.FilterFloatSeries(condition)
		for i := range df.Columns {
			df.Columns[i].RemoveIndexes(indexes)
		}
	}
}

func (df *DataFrame) FilterString(seriesName string, condition func(string) bool) {
	// Verify if series exists
	series := df.GetColumnByName(seriesName)
	if series.Name == "" {
		fmt.Println("Column not found")
	} else if series.DataType != "string" {
		fmt.Println("Not a string")
	} else {
		indexes := series.FilterStringSeries(condition)
		for i := range df.Columns {
			df.Columns[i].RemoveIndexes(indexes)
		}
	}
}

func (df *DataFrame) ReplaceWholeWord(column, old, new string) {
	name := df.GetColumnByName(column)
	name.ReplaceWholeWord(old, new)
}

func (df *DataFrame) Replace(column, old, new string) {
	name := df.GetColumnByName(column)
	name.Replace(old, new)
}

func (df *DataFrame) DropByIndex(index ...int) {
	var newSeries []Series
	for i, series := range df.Columns {
		if !ArrayContainsInteger(index, i) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
}

func (df *DataFrame) DropByName(name ...string) {
	var newSeries []Series
	for _, series := range df.Columns {
		if !ArrayContainsString(name, series.Name) {
			newSeries = append(newSeries, series)
		}
	}
	df.Columns = newSeries
}
