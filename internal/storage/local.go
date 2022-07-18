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
	metrics     map[string]gauge
	pollCounter map[string]counter
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (s *LocalStorage) defaultStorage() *LocalStorage {
	st := new(LocalStorage)
	st.metrics = make(map[string]gauge)
	st.pollCounter = make(map[string]counter)
	st.pollCounter["PollCounter"] = 0
	st.metrics["Alloc"] = 0
	st.metrics["BuckHashSys"] = 0
	st.metrics["Frees"] = 0
	st.metrics["GCCPUFraction"] = 0
	st.metrics["GCSys"] = 0
	st.metrics["HeapAlloc"] = 0
	st.metrics["HeapIdle"] = 0
	st.metrics["HeapInuse"] = 0
	st.metrics["HeapObjects"] = 0
	st.metrics["HeapReleased"] = 0
	st.metrics["HeapSys"] = 0
	st.metrics["LastGC"] = 0
	st.metrics["Lookups"] = 0
	st.metrics["MCacheInuse"] = 0
	st.metrics["MCacheSys"] = 0
	st.metrics["MSpanInuse"] = 0
	st.metrics["MSpanSys"] = 0
	st.metrics["Mallocs"] = 0
	st.metrics["NextGC"] = 0
	st.metrics["NumForcedGC"] = 0
	st.metrics["NumGC"] = 0
	st.metrics["OtherSys"] = 0
	st.metrics["PauseTotalNs"] = 0
	st.metrics["StackInuse"] = 0
	st.metrics["StackSys"] = 0
	st.metrics["Sys"] = 0
	st.metrics["TotalAlloc"] = 0
	st.metrics["RandomValue"] = 0
	return st
}

func (s *LocalStorage) Store(metricName, metricValue string, isCounter bool) error {
	if isCounter {
		v, _ := strconv.Atoi(metricValue)
		if s.pollCounter["PollCounter"] > 0 {
			s.pollCounter["PollCounter"] += counter(v)
		} else {
			s.pollCounter["PollCounter"] = counter(v)
		}
	} else {
		for key, value := range s.metrics {
			if metricName == key {
				v, _ := strconv.ParseFloat(metricValue, 64)
				s.metrics[key] = v
			} else {
				logrus.Debugf("value of key: %s not updated, value = %v", key, value)
			}
		}
	}
	return nil
}

func (s *LocalStorage) Output() error {
	bytesMemStats, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	bytesPollCounter, err := json.Marshal(s.pollCounter)
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
		return strconv.FormatInt(s.pollCounter["PollCounter"], 10), nil
	} else {
		for key, value := range s.metrics {
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
	err := enc.Encode(s.metrics)
	_ = enc.Encode(s.pollCounter)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
