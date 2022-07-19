package storage

type Storage struct {
	StorageType string
}

type Repositories interface {
	Store(metricName, metricValue string, isCounter bool)
	Output() error
	Marshal() ([]byte, error)
	SingleMetric(metricName string, isCounter bool) (string, error)
	Current() map[string]interface{}
}

func NewStorage(s *Storage) (Repositories, error) {
	switch s.StorageType {
	case "local":
		return NewLocalStorage().DefaultStorage(), nil
	default:
		return NewLocalStorage().DefaultStorage(), nil
	}
}
