package agent

import (
	"context"
	"fmt"
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

func (c counter) String() string {
	return strconv.Itoa(int(c))
}

type metricsMap map[string]gauge

type customMetrics struct {
	memstats  metricsMap
	pollCount map[string]counter
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
	client := http.Client{Timeout: 20 * time.Second,
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
		pollInterval:   viper.GetInt("agent.pollInterval"),
		reportInterval: viper.GetInt("agent.reportInterval")}
}

func urlGenerator(conf config, m map[string]gauge) (urls []string) {
	server := "http://" + conf.server + ":" + conf.port + "/update/"
	for key, value := range m {
		url := ""
		url = server + strings.Split(reflect.TypeOf(value).String(), ".")[1] + "/" + key + "/" + value.String()
		urls = append(urls, url)
	}
	return urls
}

func urlGeneratorCounter(conf config, m map[string]counter) (url string) {
	server := "http://" + conf.server + ":" + conf.port + "/update/"
	for key, value := range m {
		url = server + strings.Split(reflect.TypeOf(value).String(), ".")[1] + "/" + key + "/" + value.String()
	}
	return url
}

func (a *agent) Run() error {
	ctx := context.Background()
	c := newConfig()
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)
	metrics := new(customMetrics)
	cMetrics := make(map[string]gauge, 30)
	pollCount := make(map[string]counter, 1)
	metrics.memstats = cMetrics
	metrics.pollCount = pollCount
	metrics.newCustomMetrics()
	go func() {
		for {
			t := time.NewTicker(time.Duration(c.pollInterval) * time.Second)
			select {
			case <-t.C:
				metrics.newCustomMetrics()
			case <-ctx.Done():
				t.Stop()
			}
		}
	}()

	for {
		urls := urlGenerator(*c, metrics.memstats)
		urls = append(urls, urlGeneratorCounter(*c, metrics.pollCount))
		fmt.Println(urls)
		select {
		case <-ticker.C:

		case <-ctx.Done():
			return nil
		}
		for _, url := range urls {
			req, err := a.serv.NewRequestWithContext(ctx, http.MethodPost, url, nil, nil)
			if err != nil {
				logrus.Printf("error making request: %s\n%v", err.Error(), req)
			}
			resp, err := a.client.Do(req)
			if err != nil {
				logrus.Printf("error doing request: %s\n", err.Error())
				return err
			}
			if resp.StatusCode != http.StatusOK {
				err := resp.Body.Close()
				if err != nil {
					return err
				}
			}
		}
		urls = urls[:0]
	}

}

func (m *customMetrics) newCustomMetrics() {
	memStats := new(runtime.MemStats)
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
	m.pollCount["PollCount"]++
	m.memstats["RandomValue"] = gauge(rand.Float64())
}
