package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"pipeline/orm"
)

type CSVReader struct {
	path string
}

func New(path string) *CSVReader {
	return &CSVReader{path: path}
}

func (r *CSVReader) Read() ([]models.Product, error) {
	f, err := os.Open(r.path)
	if err != nil {
		return nil, fmt.Errorf("reader: не удалось открыть %q: %w", r.path, err)
	}
	defer f.Close()

	cr := csv.NewReader(f)
	cr.LazyQuotes = true
	cr.FieldsPerRecord = -1

	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("reader: ошибка чтения заголовка: %w", err)
	}

	var products []models.Product
	line := 1
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("reader: строка %d пропущена: %v\n", line, err)
			line++
			continue
		}
		if len(row) < 2 {
			line++
			continue
		}
		products = append(products, models.Product{
			ID:          row[0],
			Description: row[1],
		})
		line++
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("reader: файл %q пустой или нечитаемый", r.path)
	}

	return products, nil
}
