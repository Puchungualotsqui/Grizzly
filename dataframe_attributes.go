package grizzly

func (df *DataFrame) GetLength() int {
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
