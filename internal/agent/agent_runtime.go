package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/kosimovsky/tricMe/internal/repo/runtimemetrics"
)

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
		typeOfValue := func(c interface{}) string {
			switch c.(type) {
			case float64:
				return "gauge"
			case int64:
				return "counter"
			default:
				return "UnknownType"
			}
		}(value)
		url = server + typeOfValue + "/" + key + "/" + gaugeToString(value)
		urls = append(urls, url)
	}
	return urls
}

func urlGeneratorCounter(conf config, m map[string]counter) (url string) {
	server := "http://" + conf.server + ":" + conf.port + "/update/"
	for key, value := range m {
		typeOfValue := func(c interface{}) string {
			switch c.(type) {
			case float64:
				return "gauge"
			case int64:
				return "counter"
			default:
				return "UnknownType"
			}
		}(value)

		url = server + typeOfValue + "/" + key + "/" + counterToString(value)
	}
	return url
}

func (a *agent) Run() error {
	ctx := context.Background()
	c := newConfig()
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)
	metrics := runtimemetrics.NewRuntimeMetrics()
	metrics.GenerateMetrics()
	go func() {
		for {
			t := time.NewTicker(time.Duration(c.pollInterval) * time.Second)
			select {
			case <-t.C:
				metrics.GenerateMetrics()
			case <-ctx.Done():
				t.Stop()
			}
		}
	}()

	for {
		urls := urlGenerator(*c, metrics.Memstats)
		urls = append(urls, urlGeneratorCounter(*c, metrics.PollCount))
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
	}

}

func (a *agent) RunWithSerialized() error {
	ctx := context.Background()
	c := newConfig()
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)
	metrics := runtimemetrics.SerializedMetrics()
	metrics.GenerateMetrics()
	go func() {
		for {
			t := time.NewTicker(time.Duration(c.pollInterval) * time.Second)
			select {
			case <-t.C:
				metrics.GenerateMetrics()
			case <-ctx.Done():
				t.Stop()
			}
		}
	}()

	for {
		select {
		case <-ticker.C:

		case <-ctx.Done():
			return nil
		}
		for _, metric := range metrics.MetricsArray {
			url := "http://" + c.server + ":" + c.port + "/update/"
			reqBody, err := json.Marshal(metric)
			if err != nil {
				return err
			}
			bodyReader := bytes.NewReader(reqBody)

			req, err := a.serv.NewRequestWithContext(ctx, http.MethodPost, url, nil, bodyReader)
			if err != nil {
				logrus.Printf("error making request: %s\n%v", err.Error(), req)
				return err
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
	}
}

type gauge = float64

func gaugeToString(g gauge) string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

type counter = int64

func counterToString(c counter) string {
	return strconv.Itoa(int(c))
}

func (a *agent) Stop() {
	a.client.CloseIdleConnections()
}
