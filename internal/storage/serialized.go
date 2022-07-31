package storage

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kosimovsky/tricMe"
	"github.com/sirupsen/logrus"
	"strings"
)

type metricsMap struct {
	MetricsMap map[string]tricMe.Metrics
}

func NewMetricsMap() *metricsMap {
	m := make(map[string]tricMe.Metrics)
	return &metricsMap{MetricsMap: m}
}

func (m *metricsMap) Store(metric tricMe.Metrics) {
	key := generateKeyHash(metric.ID, metric.MType)
	found := false

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
}

func (m *metricsMap) SingleMetric(id, mType string) (*tricMe.Metrics, error) {
	key := generateKeyHash(id, mType)
	if value, ok := m.MetricsMap[key]; ok {
		return &value, nil
	}
	err := errors.New(fmt.Sprintf("metric %s of type %s not found", id, mType))
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
func (m *metricsMap) Output() error {
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

func (m *metricsMap) Marshal() ([]byte, error) {
	// TODO implement
	return nil, nil
}

func (m *metricsMap) Current() map[string]interface{} {
	currentMetrics := map[string]interface{}{}
	for key, value := range m.MetricsMap {
		currentMetrics[key] = value
	}
	return currentMetrics
}
