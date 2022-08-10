package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kosimovsky/tricMe/config"
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
	transport := &http.Transport{
		DialContext: (&net.Dialer{Timeout: 30 * time.Second,
			KeepAlive: 30 * time.Second}).DialContext,
		MaxIdleConns:    30,
		IdleConnTimeout: 60 * time.Second,
	}
	client := http.Client{Timeout: 5 * time.Second,
		Transport: transport}
	return &agent{client: &client,
		serv: serv}
}

func (a *agent) Run() error {
	ctx := context.Background()
	c := config.AgentConfig()
	ticker := time.NewTicker(c.Report)
	metrics := runtimemetrics.SerializedMetrics()
	metrics.GenerateMetrics()
	go func() {
		for {
			t := time.NewTicker(c.Poll)
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
			err := fmt.Errorf("timeout exceeded")
			return err
		}
		for _, metric := range metrics.MetricsArray {
			url := fmt.Sprintf("http://%s/update/", c.Address)
			reqBody, err := json.Marshal(metric)
			if err != nil {
				return err
			}
			bodyReader := bytes.NewReader(reqBody)

			req, err := a.serv.NewRequestWithContext(ctx, http.MethodPost, url, nil, bodyReader)
			if err != nil {
				logrus.Errorf("error making request: %s\n%v", err.Error(), req)
				return err
			}
			resp, err := a.client.Do(req)
			if err != nil {
				logrus.Errorf("error doing request: %s\n", err.Error())
			} else if resp.StatusCode != http.StatusOK {
				err = resp.Body.Close()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (a *agent) Stop() {
	a.client.CloseIdleConnections()
}
