package service

import (
	"context"
	"io"
	"net/http"
)

type Metrics interface {
	GenerateMetrics()
}

type service struct {
	metrics Metrics
}

func New(metrics Metrics) *service {
	return &service{metrics: metrics}
}

func (s *service) NewRequestWithContext(ctx context.Context, method, url string, headers *http.Header, body io.Reader) (*http.Request, error) {

	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		r.Header = *headers
	}
	r.Header.Set("Content-Type", "text/plain")
	r.Header.Add("Accept", "text/plain")
	return r, nil
}
