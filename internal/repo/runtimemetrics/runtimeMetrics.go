package runtimemetrics

import (
	"fmt"
	"reflect"
	"runtime"
)

type RTMetrics struct {
	memstats *runtime.MemStats
}

func New() *RTMetrics {
	m := new(runtime.MemStats)
	return &RTMetrics{memstats: m}
}

func (m *RTMetrics) GetMetricsName() string {
	return fmt.Sprintf("It is %v metrics", reflect.TypeOf(m.memstats))

}
