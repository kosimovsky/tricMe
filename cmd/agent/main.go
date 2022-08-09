package main

import (
	"fmt"
	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/agent"
	"github.com/kosimovsky/tricMe/internal/logger"
	"github.com/kosimovsky/tricMe/internal/repo"
	"github.com/kosimovsky/tricMe/internal/service"
)

func main() {
	if err := config.InitAgentConfig(); err != nil {
		_ = fmt.Errorf("error while reading config file %v", err.Error())
	}
	c := config.AgentConfig()
	agentLog := logger.NewLogger()
	agentLog.Default(c.Logfile)
	m, _ := repo.NewMiner(&repo.Source{Resources: c.MetricsType})
	serv := service.New(m)
	newAgent := agent.NewAgent(serv)

	defer newAgent.Stop()
	if err := newAgent.RunWithSerialized(); err != nil {
		agentLog.Errorf("error while running agent: %s", err.Error())
	}
}
