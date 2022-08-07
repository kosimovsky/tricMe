package storage

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	tricme "github.com/kosimovsky/tricMe"
	"github.com/sirupsen/logrus"
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
		logrus.Printf("got new metric %s of Type %s", metric.ID, metric.MType)
	}
	c := ReadConfig()
	if c.StoreInterval == 0 {
		err := m.Keep(c.Filename)
		if err != nil {
			logrus.Printf("error storing metrics: %s", err.Error())
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
func (m *metrics) Output() error {
	data, err := json.Marshal(m.MetricsMap)
	if err != nil {
		return err
	}

	output := map[string]interface{}{}
	err = json.Unmarshal(data, &output)
	if err != nil {
		return err
	}
	final, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return err
	}
	println(string(final))
	return nil
}

func (m *metrics) CurrentValues() map[string]interface{} {
	currentMetrics := map[string]interface{}{}
	for key, value := range m.MetricsMap {
		currentMetrics[key] = value
	}
	return currentMetrics
}
