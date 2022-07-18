package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

type gauge = float64
type counter = int64

type LocalStorage struct {
	gaugeMetrics   map[string]gauge
	counterMetrics map[string]counter
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (s *LocalStorage) DefaultStorage() *LocalStorage {
	st := new(LocalStorage)
	st.gaugeMetrics = make(map[string]gauge)
	st.counterMetrics = make(map[string]counter)
	st.counterMetrics["PollCounter"] = 0
	st.gaugeMetrics["Alloc"] = 0
	st.gaugeMetrics["BuckHashSys"] = 0
	st.gaugeMetrics["Frees"] = 0
	st.gaugeMetrics["GCCPUFraction"] = 0
	st.gaugeMetrics["GCSys"] = 0
	st.gaugeMetrics["HeapAlloc"] = 0
	st.gaugeMetrics["HeapIdle"] = 0
	st.gaugeMetrics["HeapInuse"] = 0
	st.gaugeMetrics["HeapObjects"] = 0
	st.gaugeMetrics["HeapReleased"] = 0
	st.gaugeMetrics["HeapSys"] = 0
	st.gaugeMetrics["LastGC"] = 0
	st.gaugeMetrics["Lookups"] = 0
	st.gaugeMetrics["MCacheInuse"] = 0
	st.gaugeMetrics["MCacheSys"] = 0
	st.gaugeMetrics["MSpanInuse"] = 0
	st.gaugeMetrics["MSpanSys"] = 0
	st.gaugeMetrics["Mallocs"] = 0
	st.gaugeMetrics["NextGC"] = 0
	st.gaugeMetrics["NumForcedGC"] = 0
	st.gaugeMetrics["NumGC"] = 0
	st.gaugeMetrics["OtherSys"] = 0
	st.gaugeMetrics["PauseTotalNs"] = 0
	st.gaugeMetrics["StackInuse"] = 0
	st.gaugeMetrics["StackSys"] = 0
	st.gaugeMetrics["Sys"] = 0
	st.gaugeMetrics["TotalAlloc"] = 0
	st.gaugeMetrics["RandomValue"] = 0
	return st
}

func (s *LocalStorage) Store(metricName, metricValue string, isCounter bool) {
	found := false
	if isCounter {
		v, _ := strconv.ParseInt(metricValue, 10, 64)
		for key, value := range s.counterMetrics {
			if metricName == key {
				found = true
				if value == 0 {
					s.counterMetrics[metricName] = v
				} else {
					s.counterMetrics[metricName] += v
				}
			}
		}
		if !found {
			s.counterMetrics[metricName] = v
			logrus.Printf("got new COUNTER metric %s with value %s", metricName, metricValue)
		}
	} else {
		v, _ := strconv.ParseFloat(metricValue, 64)
		for key, _ := range s.gaugeMetrics {
			if metricName == key {
				found = true
				s.gaugeMetrics[key] = v
			}
		}
		if !found {
			v, _ := strconv.ParseFloat(metricValue, 64)
			s.gaugeMetrics[metricName] = v
			logrus.Printf("got new GAUGE metric %s with value %s", metricName, metricValue)
		}
	}
}

func (s *LocalStorage) Output() error {
	bytesMemStats, err := json.Marshal(s.gaugeMetrics)
	if err != nil {
		return err
	}
	bytesPollCounter, err := json.Marshal(s.counterMetrics)
	if err != nil {
		return err
	}

	mergedOutput := map[string]interface{}{}

	err = json.Unmarshal(bytesMemStats, &mergedOutput)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytesPollCounter, &mergedOutput)
	if err != nil {
		return err
	}
	finalOut, err := json.MarshalIndent(mergedOutput, "", " ")
	if err != nil {
		return err
	}
	println(string(finalOut))
	return nil
}

func (s *LocalStorage) SingleMetric(metricName string, isCounter bool) (string, error) {
	if isCounter {
		for key, value := range s.counterMetrics {
			if metricName == key {
				return strconv.FormatInt(value, 10), nil
			}
		}
	} else {
		for key, value := range s.gaugeMetrics {
			if metricName == key {
				return strconv.FormatFloat(value, 'f', -1, 64), nil
			}
		}
	}
	return "None", fmt.Errorf("there is no such metric %s", metricName)
}

func (s *LocalStorage) Marshal() ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&s)
	_ = enc.Encode(s.counterMetrics)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (s *LocalStorage) Current() map[string]interface{} {
	currentMetrics := map[string]interface{}{}
	for key, value := range s.gaugeMetrics {
		currentMetrics[key] = value
	}
	for key, value := range s.counterMetrics {
		currentMetrics[key] = value
	}
	return currentMetrics
}
