package repo

import "github.com/kosimovsky/tricMe/internal/repo/runtimemetrics"

type Source struct {
	Resources string
}

type Miner interface {
	GenerateMetrics()
}

func NewMiner(s *Source) (Miner, error) {

	if s.Resources == "memStat" {
		return runtimemetrics.NewCustomMetrics(), nil
	}
	return nil, nil
}
