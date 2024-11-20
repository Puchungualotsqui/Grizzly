package grizzly

import "fmt"

func (df *DataFrame) GetLength() int {
	if len(df.Columns) == 0 {
		return 0
	}
	series := df.Columns[0]
	return series.GetLength()
}

func (df *DataFrame) GetNumberOfColumns() int {
	return len(df.Columns)
}

func (df *DataFrame) GetColumnNames() []string {
	var names []string
	for _, series := range df.Columns {
		names = append(names, series.Name)
	}
	return names
}

func (df *DataFrame) GetShape() [2]int {
	var shape [2]int
	shape[0] = len(df.Columns)
	shape[1] = df.GetLength()
	return shape
}

func (df *DataFrame) ContainsColumn(name string) bool {
	names := df.GetColumnNames()
	return ArrayContainsString(names, name)
}

func (df *DataFrame) GetColumnByName(name string) *Series {
	for i, series := range df.Columns {
		if series.Name == name {
			return &df.Columns[i]
		}
	}
	panic(fmt.Sprintf("%s not found", name))
}

func (df *DataFrame) GetColumnByIndex(index int) *Series {
	return &df.Columns[index]
}
