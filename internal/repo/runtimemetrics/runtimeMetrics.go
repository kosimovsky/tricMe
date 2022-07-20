package runtimemetrics

import (
	"math/rand"
	"runtime"
)

type gauge = float64

type counter = int64

type CustomMetrics struct {
	Memstats  map[string]gauge
	PollCount map[string]counter
}

func NewCustomMetrics() *CustomMetrics {
	m := new(CustomMetrics)
	cMetrics := make(map[string]gauge, 30)
	pollCount := make(map[string]counter, 1)
	m.Memstats = cMetrics
	m.PollCount = pollCount
	return m
}

func (m *CustomMetrics) GenerateMetrics() {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	m.Memstats["Alloc"] = gauge(memStats.Alloc)
	m.Memstats["BuckHashSys"] = gauge(memStats.BuckHashSys)
	m.Memstats["Frees"] = gauge(memStats.Frees)
	m.Memstats["GCCPUFraction"] = gauge(memStats.GCCPUFraction)
	m.Memstats["GCSys"] = gauge(memStats.GCSys)
	m.Memstats["HeapAlloc"] = gauge(memStats.HeapAlloc)
	m.Memstats["HeapIdle"] = gauge(memStats.HeapIdle)
	m.Memstats["HeapInuse"] = gauge(memStats.HeapInuse)
	m.Memstats["HeapObjects"] = gauge(memStats.HeapObjects)
	m.Memstats["HeapReleased"] = gauge(memStats.HeapReleased)
	m.Memstats["HeapSys"] = gauge(memStats.HeapSys)
	m.Memstats["LastGC"] = gauge(memStats.LastGC)
	m.Memstats["Lookups"] = gauge(memStats.Lookups)
	m.Memstats["MCacheInuse"] = gauge(memStats.MCacheInuse)
	m.Memstats["MCacheSys"] = gauge(memStats.MCacheSys)
	m.Memstats["MSpanInuse"] = gauge(memStats.MSpanInuse)
	m.Memstats["MSpanSys"] = gauge(memStats.MSpanSys)
	m.Memstats["Mallocs"] = gauge(memStats.Mallocs)
	m.Memstats["NextGC"] = gauge(memStats.NextGC)
	m.Memstats["NumForcedGC"] = gauge(memStats.NumForcedGC)
	m.Memstats["NumGC"] = gauge(memStats.NumGC)
	m.Memstats["OtherSys"] = gauge(memStats.OtherSys)
	m.Memstats["PauseTotalNs"] = gauge(memStats.PauseTotalNs)
	m.Memstats["StackInuse"] = gauge(memStats.StackInuse)
	m.Memstats["StackSys"] = gauge(memStats.StackSys)
	m.Memstats["Sys"] = gauge(memStats.Sys)
	m.Memstats["TotalAlloc"] = gauge(memStats.TotalAlloc)
	m.PollCount["PollCount"]++
	m.Memstats["RandomValue"] = gauge(rand.Float64())
}
