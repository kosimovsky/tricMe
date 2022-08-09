package storage

import (
	tricme "github.com/kosimovsky/tricMe"
	"github.com/spf13/viper"
)

type Storage struct {
	StorageType string
}

type Storekeeper interface {
	Store(metrics tricme.Metrics)
	Output() error
	Keep(file string) error
	Restore(filename string, flag bool) error
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
	case "test":
		return TestMetrics(), nil
	default:
		return NewMetricsMap(), nil
	}
}
