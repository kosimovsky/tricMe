package agent

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type gauge float64

func (g gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

type counter int64

type metricsMap map[string]gauge

type customMetrics struct {
	memstats  metricsMap
	PollCount counter
}

type Sender interface {
	NewRequestWithContext(ctx context.Context, method, url string, headers *http.Header, body io.Reader) (*http.Request, error)
}

type agent struct {
	client *http.Client
	serv   Sender
}

func NewAgent(serv Sender) *agent {
	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	client := http.Client{Timeout: 2 * time.Second,
		Transport: transport}
	return &agent{client: &client,
		serv: serv}
}

type config struct {
	server         string
	port           string
	pollInterval   int
	reportInterval int
}

func newConfig() *config {
	return &config{server: viper.GetString("server.address"),
		port:           viper.GetString("server.port"),
		pollInterval:   viper.GetInt("agent.pollIntervall"),
		reportInterval: viper.GetInt("agent.reportInterval")}
}

func urlGenerator(conf config, m map[string]gauge) (urls []string) {
	var server string
	server = "http://" + conf.server + ":" + conf.port + "/update/"
	url := server
	for key, value := range m {
		url += url + strings.Split(reflect.TypeOf(value).String(), ".")[1] + "/" + key + value.String()
		urls = append(urls, url)
	}
	return urls
}

func (a *agent) Run() error {
	ctx := context.Background()
	c := newConfig()
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)

	metrics := new(customMetrics)
	cMetrics := make(map[string]gauge, 30)

	metrics.memstats = cMetrics

	metrics.newCustomMetrics(ctx, 2)

	urls := urlGenerator(*c, metrics.memstats)

	for _, url := range urls {
		req, err := a.serv.NewRequestWithContext(ctx, http.MethodPost, url, nil, nil)
		if err != nil {
			logrus.Printf("error making request: %s\n%v", err.Error(), req)
		}
		select {
		case <-ticker.C:
			_, err := a.client.Do(req)
			if err != nil {
				logrus.Printf("error doing request: %s\n", err.Error())
			}
		case <-ctx.Done():
			return nil
		}

	}
	return nil
}

func (m *customMetrics) newCustomMetrics(ctx context.Context, pollInterval int) customMetrics {
	memStats := new(runtime.MemStats)
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	select {
	case <-ticker.C:
		runtime.ReadMemStats(memStats)
		m.memstats["Alloc"] = gauge(memStats.Alloc)
		m.memstats["BuckHashSys"] = gauge(memStats.BuckHashSys)
		m.memstats["Frees"] = gauge(memStats.Frees)
		m.memstats["GCCPUFraction"] = gauge(memStats.GCCPUFraction)
		m.memstats["GCSys"] = gauge(memStats.GCSys)
		m.memstats["HeapAlloc"] = gauge(memStats.HeapAlloc)
		m.memstats["HeapIdle"] = gauge(memStats.HeapIdle)
		m.memstats["HeapInuse"] = gauge(memStats.HeapInuse)
		m.memstats["HeapObjects"] = gauge(memStats.HeapObjects)
		m.memstats["HeapReleased"] = gauge(memStats.HeapReleased)
		m.memstats["HeapSys"] = gauge(memStats.HeapSys)
		m.memstats["LastGC"] = gauge(memStats.LastGC)
		m.memstats["Lookups"] = gauge(memStats.Lookups)
		m.memstats["MCacheInuse"] = gauge(memStats.MCacheInuse)
		m.memstats["MCacheSys"] = gauge(memStats.MCacheSys)
		m.memstats["MSpanInuse"] = gauge(memStats.MSpanInuse)
		m.memstats["MSpanSys"] = gauge(memStats.MSpanSys)
		m.memstats["Mallocs"] = gauge(memStats.Mallocs)
		m.memstats["NextGC"] = gauge(memStats.NextGC)
		m.memstats["NumForcedGC"] = gauge(memStats.NumForcedGC)
		m.memstats["NumGC"] = gauge(memStats.NumGC)
		m.memstats["OtherSys"] = gauge(memStats.OtherSys)
		m.memstats["PauseTotalNs"] = gauge(memStats.PauseTotalNs)
		m.memstats["StackInuse"] = gauge(memStats.StackInuse)
		m.memstats["StackSys"] = gauge(memStats.StackSys)
		m.memstats["Sys"] = gauge(memStats.Sys)
		m.memstats["TotalAlloc"] = gauge(memStats.TotalAlloc)
		m.PollCount++
		m.memstats["RandomValue"] = gauge(rand.Float64())
		return *m
	case <-ctx.Done():
		return *m.defaultCM()
	}
}

func (m *customMetrics) defaultCM() *customMetrics {
	return nil
}
