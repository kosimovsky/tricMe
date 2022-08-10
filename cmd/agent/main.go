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
		fmt.Printf("error while reading config file %s", err.Error())
	}
	c := config.AgentConfig()
	agentLog := logger.NewLogger()
	agentLog.Default(c.Logfile)
	m, err := repo.NewMiner(&repo.Source{Resources: c.MetricsType})
	if err != nil {
		agentLog.Error(err)
	}
	serv := service.New(m)
	newAgent := agent.NewAgent(serv)

	defer newAgent.Stop()
	if err = newAgent.Run(); err != nil {
		agentLog.Errorf("error while running agent: %s", err.Error())
	}
}
