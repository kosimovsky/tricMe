package storage

type Storage struct {
	StorageType string
}

type Repositories interface {
	Store(string) error
	Output() error
}

func NewStorage(s *Storage) (Repositories, error) {
	switch s.StorageType {
	case "local":
		return NewLocalStorage().defaultStorage(), nil
	default:
		return NewLocalStorage().defaultStorage(), nil
	}
}
