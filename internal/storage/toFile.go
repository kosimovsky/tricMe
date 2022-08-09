package storage

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

type Storer interface {
	WriteMetric(metric *metrics) error
	Close() error
}

type Restorer interface {
	ReadMetric() (*metrics, error)
	Close() error
}

type storer struct {
	file    *os.File
	encoder *json.Encoder
}

func newStorer(filename string) (*storer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &storer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

type restorer struct {
	file    *os.File
	decoder *json.Decoder
}

func newRestorer(filename string) (*restorer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &restorer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *storer) WriteMetric(metrics *metrics) error {
	return p.encoder.Encode(&metrics)
}

func (p *storer) Close() error {
	return p.file.Close()
}

func (r *restorer) ReadMetric() (*metrics, error) {
	m := &metrics{}
	if err := r.decoder.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *restorer) Close() error {
	return r.file.Close()
}

func (m *metrics) Keep(filename string) error {
	s, err := newStorer(filename)
	if err != nil {
		return err
	}
	defer s.Close()
	m.mx.Lock()
	defer m.mx.Unlock()
	err = s.WriteMetric(m)
	if err != nil {
		logrus.Error(err.Error())
	}
	return nil
}

func (m *metrics) Restore(filename string, flag bool) error {
	if flag {
		r, err := newRestorer(filename)
		if err != nil {
			return err
		}
		metric, err := r.ReadMetric()
		if err != nil {
			return err
		}
		m.MetricsMap = metric.MetricsMap
	}
	return nil
}
