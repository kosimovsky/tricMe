package storage

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	tricme "github.com/kosimovsky/tricMe"
	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/log"
)

type metrics struct {
	mx         sync.RWMutex
	MetricsMap map[string]tricme.Metrics
}

func NewMetricsMap() *metrics {
	m := make(map[string]tricme.Metrics)
	return &metrics{MetricsMap: m}
}

func (m *metrics) Store(metric tricme.Metrics) {
	c := config.NewServerConfig()
	logger, err := log.NewLogger(c.Logfile)
	if err != nil {
		fmt.Printf("error occurs while logger initialization: %s", err.Error())
		return
	}
	defer logger.Close()
	key := generateKeyHash(metric.ID, metric.MType)
	found := false

	m.mx.Lock()
	defer m.mx.Unlock()
	for k, v := range m.MetricsMap {
		if k == key {
			found = true
			if metric.Delta != nil {
				*v.Delta += *metric.Delta
			}
			if metric.Value != nil {
				*v.Value = *metric.Value
			}
		}
	}
	if !found {
		m.MetricsMap[key] = metric
		logger.Printf("got new metric %s of Type %s", metric.ID, metric.MType)
	}

	if c.StoreInterval == 0 {
		err = m.keepDirectly(c.Filename)
		if err != nil {
			logger.Printf("error storing metrics: %s", err.Error())
		}
	}

}

func (m *metrics) SingleMetric(id, mType string) (*tricme.Metrics, error) {
	key := generateKeyHash(id, mType)
	m.mx.RLock()
	defer m.mx.RUnlock()
	if value, ok := m.MetricsMap[key]; ok {
		return &value, nil
	}
	err := fmt.Errorf("metric %s of type %s not found", id, mType)
	return nil, err
}

// generateKeyHash generates hash from ID and MType of tricMe.Metrics struct for key to store Metric in map
func generateKeyHash(id, mType string) string {
	hash := sha1.New()
	var sb strings.Builder
	sb.WriteString(id)
	sb.WriteString(mType)
	hash.Write([]byte(sb.String()))
	return fmt.Sprintf("%x", hash.Sum([]byte(nil)))
}

// Output is for debugging server. To start output every 5 seconds set server.debug to True in config
func (m *metrics) Output(logger *log.Logger) {
	ticker := time.NewTicker(5 * time.Second)
	ctx := context.Background()
	for {
		select {
		case <-ticker.C:
			data, err := json.Marshal(m.MetricsMap)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
			output := map[string]interface{}{}
			err = json.Unmarshal(data, &output)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
			final, err := json.MarshalIndent(output, "", " ")
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
			println(string(final))
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (m *metrics) CurrentValues() map[string]interface{} {
	currentMetrics := map[string]interface{}{}
	for key, value := range m.MetricsMap {
		currentMetrics[key] = value
	}
	return currentMetrics
}

func TestMetrics() *metrics {
	mMap := new(metrics)
	rv := 0.4246374970712657
	pollCount := int64(15)
	totallAlloc := float64(2794104)
	mallocs := float64(23543)
	alloc := float64(2794104)

	tMap := map[string]tricme.Metrics{
		"03e66670c529012c38b396b1872d680ead69f624": {
			ID:    "RandomValue",
			MType: "gauge",
			Value: &rv,
		},
		"076f524e410c90c19cd689e8f00e598556cf4468": {
			Delta: &pollCount,
			ID:    "PollCount",
			MType: "counter",
		},
		"330f5a124e1cc3957b90ad0f3add29fa1ba58a1c": {
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: &totallAlloc,
		},
		"d050d8574c1fc8794832f6e2eca94604ba103529": {
			ID:    "Mallocs",
			MType: "gauge",
			Value: &mallocs,
		},
		"f4bf2f6d4d4e2c54f0e23bc8268dfc7531e37653": {
			ID:    "Alloc",
			MType: "gauge",
			Value: &alloc,
		},
	}
	mMap.MetricsMap = tMap
	return mMap
}
