package runtimemetrics

import (
	"github.com/kosimovsky/tricMe"
)

type serializedMetrics struct {
	runtimeMetrics *RuntimeMetrics
	MetricsArray   []tricMe.Metrics
}

func SerializedMetrics() *serializedMetrics {
	arr := new([]tricMe.Metrics)
	return &serializedMetrics{
		runtimeMetrics: NewRuntimeMetrics(),
		MetricsArray:   *arr,
	}
}

func (c *serializedMetrics) GenerateMetrics() {
	c.runtimeMetrics.GenerateMetrics()
	for key, value := range c.runtimeMetrics.Memstats {
		tmp := value
		metric := new(tricMe.Metrics)
		metric.ID = key
		metric.MType = "gauge"
		metric.Value = &tmp
		c.MetricsArray = append(c.MetricsArray, *metric)
	}
	for key, value := range c.runtimeMetrics.PollCount {
		metric := new(tricMe.Metrics)
		tmp := value
		metric.ID = key
		metric.MType = "counter"
		metric.Delta = &tmp
		c.MetricsArray = append(c.MetricsArray, *metric)
	}
}
