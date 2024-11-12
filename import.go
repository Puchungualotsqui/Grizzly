package grizzly

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func ImportCSV(filepath string) (DataFrame, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return DataFrame{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return DataFrame{}, fmt.Errorf("failed to read CSV file: %v", err)
	}

	if len(records) == 0 {
		return DataFrame{}, fmt.Errorf("CSV file is empty")
	}

	headers := records[0]
	columns := make([]Series, len(headers))

	// Initialize Series for each header
	for i, header := range headers {
		columns[i] = Series{Name: header, DataType: "string", String: []string{}}
	}

	// Populate Series with data
	for _, row := range records[1:] {
		for i, value := range row {
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				columns[i].Float = append(columns[i].Float, floatValue)
				columns[i].DataType = "float"
			} else {
				columns[i].String = append(columns[i].String, value)
			}
		}
	}

	return DataFrame{Columns: columns}, nil
}
