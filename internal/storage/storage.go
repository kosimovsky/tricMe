package storage

import (
	"github.com/spf13/viper"

	tricme "github.com/kosimovsky/tricMe"
)

type Storage struct {
	StorageType string
}

type Storekeeper interface {
	Store(metrics tricme.Metrics)
	Output() error
	Keep() error
	Restore() error
	SingleMetric(id, mType string) (*tricme.Metrics, error)
	CurrentValues() map[string]interface{}
}

func NewStorage(s *Storage) (Storekeeper, error) {
	if viper.GetString("File") == "" {
		s.StorageType = ""
	}

	switch s.StorageType {
	case "memory":
		return NewMetricsMap(), nil
	default:
		return NewMetricsMap(), nil
	}
}
