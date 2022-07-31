package runtimemetrics

import tricme "github.com/kosimovsky/tricMe"

type serializedMetrics struct {
	runtimeMetrics *RuntimeMetrics
	MetricsArray   []tricme.Metrics
}

func SerializedMetrics() *serializedMetrics {
	arr := new([]tricme.Metrics)
	return &serializedMetrics{
		runtimeMetrics: NewRuntimeMetrics(),
		MetricsArray:   *arr,
	}
}

func (c *serializedMetrics) GenerateMetrics() {
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
}
