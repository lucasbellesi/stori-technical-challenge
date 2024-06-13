package transactions

import (
	"encoding/csv"
	"os"
)

type CSVReader interface {
	Read(filePath string) ([][]string, error)
}

type DefaultCSVReader struct{}

func (r DefaultCSVReader) Read(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}
