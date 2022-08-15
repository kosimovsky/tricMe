package repo

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/kosimovsky/tricMe/internal/repo/runtimemetrics"
)

type source struct {
}

type Miner interface {
	GenerateMetrics(ctx context.Context, ticker *time.Ticker)
	NewRequestWithContext(ctx context.Context, method, url string, headers *http.Header, body io.Reader) (*http.Request, error)
}

func NewMiner(source string) (Miner, error) {
	if source == "memStat" {
		return runtimemetrics.NewMetrics(), nil
	}
	return nil, nil
}
