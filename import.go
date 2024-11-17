package grizzly

import (
	"encoding/csv"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func ImportCSV(filepath string, extra ...int) DataFrame {
	file, err := os.Open(filepath)
	if err != nil {
		panic("File was not found")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		panic("Error reading CSV headers")
	}
	numCols := len(headers)

	columns := make([]Series, numCols)
	for i, header := range headers {
		columns[i] = Series{
			Name:     header,
			DataType: "string",
			String:   make([]string, 0, 1000),
			Float:    make([]float64, 0, 1000),
		}
	}

	numGoroutines := runtime.NumCPU()
	rowChannel := make(chan [][]string, numGoroutines*2)
	var wg sync.WaitGroup

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rows := range rowChannel {
				for _, row := range rows {
					for i, value := range row {
						if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
							columns[i].Float = append(columns[i].Float, floatValue)
							columns[i].DataType = "float"
						} else {
							columns[i].String = append(columns[i].String, value)
						}
					}
				}
			}
		}()
	}

	// Get Batch Size from optional input
	var batchSize int
	if len(extra) > 0 {
		batchSize = 100
	} else {
		batchSize = extra[0]
	}
	batch := make([][]string, 0, batchSize)
	for {
		row, err := reader.Read()
		if err != nil {
			if len(batch) > 0 {
				rowChannel <- batch
			}
			break
		}
		batch = append(batch, row)
		if len(batch) >= batchSize {
			rowChannel <- batch
			batch = make([][]string, 0, batchSize)
		}
	}
	close(rowChannel)
	wg.Wait()

	return DataFrame{Columns: columns}
}

func ImportCSVOld(filepath string) DataFrame {
	file, err := os.Open(filepath)
	if err != nil {
		panic("File was not found")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if len(records) == 0 {
		return DataFrame{}
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

	return DataFrame{Columns: columns}
}
