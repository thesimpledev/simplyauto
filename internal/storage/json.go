package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const FileExtension = ".simplyauto"

type JSONStorage struct{}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{}
}

func (s *JSONStorage) Save(recording *Recording, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if filepath.Ext(path) == "" {
		path += FileExtension
	}

	recording.Finalize()

	data, err := json.MarshalIndent(recording, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal recording: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *JSONStorage) Load(path string) (*Recording, error) {
	// Check if file exists; if not and no extension, try default extension
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if filepath.Ext(path) == "" {
			path += FileExtension
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var recording Recording
	if err := json.Unmarshal(data, &recording); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recording: %w", err)
	}

	return &recording, nil
}

func (s *JSONStorage) GetDefaultExtension() string {
	return FileExtension
}
