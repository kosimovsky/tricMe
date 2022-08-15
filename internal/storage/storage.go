package storage

import (
	"time"

	tricme "github.com/kosimovsky/tricMe"
	"github.com/kosimovsky/tricMe/internal/log"
)

type storage struct {
}

type Storekeeper interface {
	Store(metrics tricme.Metrics)
	Output(logger *log.Logger)
	Keep(file string, interval time.Duration, logger *log.Logger)
	Restore(filename string, flag bool) error
	SingleMetric(id, mType string) (*tricme.Metrics, error)
	CurrentValues() map[string]interface{}
}

func NewStorage(storageType string) (Storekeeper, error) {

	switch storageType {
	case "memory":
		return NewMetricsMap(), nil
	case "test":
		return TestMetrics(), nil
	default:
		return NewMetricsMap(), nil
	}
}
