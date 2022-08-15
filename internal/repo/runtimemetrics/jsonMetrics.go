package runtimemetrics

import (
	"context"
	"io"
	"net/http"
	"time"

	tricme "github.com/kosimovsky/tricMe"
)

type metrics struct {
	runtimeMetrics *runtimeMetrics
	MetricsArray   []tricme.Metrics
}

func NewMetrics() *metrics {
	arr := new([]tricme.Metrics)
	return &metrics{
		runtimeMetrics: NewRuntimeMetrics(),
		MetricsArray:   *arr,
	}
}

func (c *metrics) GenerateMetrics(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			c.runtimeMetrics.GenerateMetrics()
			for key, value := range c.runtimeMetrics.Memstats {
				tmp := value
				metric := new(tricme.Metrics)
				metric.ID = key
				metric.MType = "gauge"
				metric.Value = &tmp
				c.MetricsArray = append(c.MetricsArray, *metric)
			}
			for key, value := range c.runtimeMetrics.PollCount {
				metric := new(tricme.Metrics)
				tmp := value
				metric.ID = key
				metric.MType = "counter"
				metric.Delta = &tmp
				c.MetricsArray = append(c.MetricsArray, *metric)
			}
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (c *metrics) NewRequestWithContext(ctx context.Context, method, url string, headers *http.Header, body io.Reader) (*http.Request, error) {

	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		r.Header = *headers
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	return r, nil
}
