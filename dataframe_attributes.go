package grizzly

import (
	"fmt"
)

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
	return arrayContainsString(names, name)
}

func (df *DataFrame) GetColumnByName(name string) (*Series, error) {
	for i, series := range df.Columns {
		if series.Name == name {
			return &df.Columns[i], nil
		}
	}
	return nil, fmt.Errorf("column %q not found", name)
}

func (df *DataFrame) GetColumnByIndex(index int) (*Series, error) {
	if index < 0 || index >= len(df.Columns) {
		return nil, fmt.Errorf("index %d is out of bounds", index)
	}
	return &df.Columns[index], nil
}

func (df *DataFrame) GetColumnDynamic(identifier any) (*Series, error) {
	var possibleName string
	var possibleIndex int
	var byIndex bool
	var err error
	var result *Series

	switch v := identifier.(type) {
	case int:
		byIndex = true
		possibleIndex = v
	default:
		byIndex = false
		possibleName, err = interfaceConvertToString(identifier)
	}
	if byIndex {
		result, err = df.GetColumnByIndex(possibleIndex)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	result, err = df.GetColumnByName(possibleName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (df *DataFrame) GetColumnType(identifier any) (string, error) {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return "", err
	}
	return series.DataType, nil
}

func (df *DataFrame) ColumnIsString(identifier any) (bool, error) {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return false, err
	}
	return series.DataType == "string", nil
}

func (df *DataFrame) ColumnIsFloat(identifier any) (bool, error) {
	series, err := df.GetColumnDynamic(identifier)
	if err != nil {
		return false, err
	}
	return series.DataType == "float", nil
}

func (df *DataFrame) GetColumnIndexByName(columnName string) (int, error) {
	for i, series := range df.Columns {
		if series.Name == columnName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column %q not found", columnName)
}
