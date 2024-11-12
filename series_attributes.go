package grizzly

func (series *Series) GetLength() int {
	if series.DataType == "float" {
		return len(series.Float)
	} else {
		return len(series.String)
	}
}
