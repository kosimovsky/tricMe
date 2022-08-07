package storage

import (
	"github.com/spf13/viper"
	"time"

	tricme "github.com/kosimovsky/tricMe"
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
	default:
		return NewMetricsMap(), nil
	}
}

type config struct {
	Address       string
	Debug         bool
	GinMode       string
	Logfile       string
	Loglevel      int
	Storage       string
	StoreInterval time.Duration
	Filename      string
	Restore       bool
}

func ReadConfig() *config {
	return &config{
		Address:       viper.GetString("Address"),
		Debug:         viper.GetBool("Debug"),
		GinMode:       viper.GetString("gonMode"),
		Logfile:       viper.GetString("Logfile"),
		Loglevel:      viper.GetInt("Loglevel"),
		Storage:       viper.GetString("Storage"),
		StoreInterval: viper.GetDuration("Interval"),
		Filename:      viper.GetString("File"),
		Restore:       viper.GetBool("Restore"),
	}
}
