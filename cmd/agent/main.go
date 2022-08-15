package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/agent"
	"github.com/kosimovsky/tricMe/internal/log"
	"github.com/kosimovsky/tricMe/internal/repo"
)

func main() {
	if err := config.InitAgentConfig(); err != nil {
		panic(fmt.Sprintf("error while reading config file %s", err.Error()))
	}
	c := config.NewAgentConfig()
	logger, err := log.NewLogger(c.Logfile)
	if err != nil {
		fmt.Printf("error occurs while logger initialization: %s", err.Error())
		os.Exit(1)
	}
	miner, err := repo.NewMiner(c.MetricsType)
	if err != nil {
		logger.Error(err)
	}
	newAgent := agent.NewAgent(miner, c)

	defer newAgent.Stop()
	if err = newAgent.Run(logger); err != nil {
		logger.Fatalf("error while running agent: %s", err.Error())
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	logger.Printf("Recieved a signal %v. Agent is shutting down...", sig)
}
