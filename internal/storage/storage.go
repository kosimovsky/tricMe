package storage

import "github.com/kosimovsky/tricMe"

type Storage struct {
	StorageType string
}

type Storekeeper interface {
	Store(metrics tricMe.Metrics)
	Output() error
	Marshal() ([]byte, error)
	SingleMetric(id, mType string) (*tricMe.Metrics, error)
	Current() map[string]interface{}
}

func NewStorage(s *Storage) (Storekeeper, error) {
	switch s.StorageType {
	case "local":
		return NewMetricsMap(), nil
	default:
		return NewMetricsMap(), nil
	}
}
