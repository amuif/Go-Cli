package main

import (
	"encoding/json"
	"os"
)

type Storage[T any] struct {
	FileName string
}

func NewStorage[T any](filName string) *Storage[T] {
	return &Storage[T]{FileName: filName}
}

func (s *Storage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "   ")
  if err != nil {
    return nil 
  }
  return os.WriteFile(s.FileName, fileData,0644)
}

func (s *Storage[T]) Load(data *T) error {
  fileData, err := os.ReadFile(s.FileName)
  if err != nil {
    return nil 
  }
  return json.Unmarshal(fileData, data)
}
