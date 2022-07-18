package storage

type Storage struct {
	StorageType string
}

type Repositories interface {
	Store(metricName, metricValue string, isCounter bool) error
	Output() error
	Marshal() ([]byte, error)
	SingleMetric(metricName string, isCounter bool) (string, error)
}

func NewStorage(s *Storage) (Repositories, error) {
	switch s.StorageType {
	case "local":
		return NewLocalStorage().defaultStorage(), nil
	default:
		return NewLocalStorage().defaultStorage(), nil
	}
}
