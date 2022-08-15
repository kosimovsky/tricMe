package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/log"
	"github.com/kosimovsky/tricMe/internal/repo"
	"github.com/kosimovsky/tricMe/internal/repo/runtimemetrics"
)

type agent struct {
	client *http.Client
	miner  repo.Miner

	pollTicker   *time.Ticker
	reportTicker *time.Ticker

	address string
}

func NewAgent(miner repo.Miner, c *config.AgentConfig) *agent {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		DialContext:     dialer.DialContext,
		MaxIdleConns:    30,
		IdleConnTimeout: 60 * time.Second,
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	return &agent{
		client:       client,
		miner:        miner,
		pollTicker:   time.NewTicker(c.Poll),
		reportTicker: time.NewTicker(c.Report),
		address:      c.Address,
	}
}

func (a *agent) Run(logger *log.Logger) error {
	ctx := context.Background()
	url := fmt.Sprintf("http://%s/update/", a.address)
	metrics := runtimemetrics.NewMetrics()

	go metrics.GenerateMetrics(ctx, a.pollTicker)

	for {
		select {
		case <-a.reportTicker.C:

		case <-ctx.Done():
			err := fmt.Errorf("timeout exceeded")
			return err
		}
		for _, metric := range metrics.MetricsArray {

			reqBody, err := json.Marshal(metric)
			if err != nil {
				return err
			}
			bodyReader := bytes.NewReader(reqBody)

			req, err := a.miner.NewRequestWithContext(ctx, http.MethodPost, url, nil, bodyReader)
			if err != nil {
				logger.Errorf("error making request: %s\n%v", err.Error(), req)
				return err
			}
			resp, err := a.client.Do(req)
			if err != nil {
				logger.Errorf("error doing request: %s\n", err.Error())
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
