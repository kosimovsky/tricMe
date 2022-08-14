package storage

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/kosimovsky/tricMe/internal/log"
)

type Storer interface {
	WriteMetric(metric *metrics) error
	Close() error
}

type Restorer interface {
	ReadMetric() (*metrics, error)
	Close() error
}

type store struct {
	file    *os.File
	encoder *json.Encoder

	mu sync.RWMutex
}

func newStore(filename string) (*store, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &store{
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

func (p *store) WriteMetric(metrics *metrics) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.encoder.Encode(&metrics)
}

func (p *store) Close() error {
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

func (m *metrics) Keep(filename string, interval time.Duration, logger *log.Logger) {
	ctx := context.Background()
	ticker := time.NewTicker(interval)
	select {
	case <-ticker.C:
		{
			s, err := newStore(filename)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
			defer s.Close()
			err = s.WriteMetric(m)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}
		}
	case <-ctx.Done():
		ticker.Stop()
	}
}

func (m *metrics) keepDirectly(filename string) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	s, err := newStore(filename)
	if err != nil {
		return err
	}
	defer s.Close()
	err = s.WriteMetric(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *metrics) Restore(filename string, flag bool) error {
	if flag {
		r, err := newRestorer(filename)
		if err != nil {
			return err
		}
		defer r.Close()
		metric, err := r.ReadMetric()
		if err != nil {
			return err
		}
		m.MetricsMap = metric.MetricsMap
	}
	return nil
}
