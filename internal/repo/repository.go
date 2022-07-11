package repo

import "github.com/kosimovsky/tricMe/internal/repo/runtimeMetrics"

type Source struct {
	Resources string
}

type Miner interface {
	GetMetricsName() string
}

func NewMiner(s *Source) (Miner, error) {

	if s.Resources == "memStat" {
		return runtimeMetrics.New(), nil
	}
	return nil, nil
}
