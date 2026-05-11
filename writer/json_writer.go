package writer

import (
	"encoding/json"
	"fmt"
	"os"

	"pipeline/orm"
)

type JSONWriter struct {
	path string
}

func New(path string) *JSONWriter {
	return &JSONWriter{path: path}
}

func (w *JSONWriter) Write(result models.PipelineResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("writer: marshal: %w", err)
	}
	if err := os.WriteFile(w.path, data, 0644); err != nil {
		return fmt.Errorf("writer: запись %q: %w", w.path, err)
	}
	return nil
}
