package main

import (
	"fmt"
	"github.com/kosimovsky/tricMe/config"
	"github.com/kosimovsky/tricMe/internal/agent"
	"github.com/kosimovsky/tricMe/internal/repo"
	"github.com/kosimovsky/tricMe/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if err := config.InitAgentConfig(); err != nil {
		_ = fmt.Errorf("error while reading config file %v", err.Error())
	}
	logfileFromConfig := viper.GetString("Logfile")
	logfile, err := os.OpenFile(logfileFromConfig, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error occured opening file : %v, %s", err, logfileFromConfig)
	}
	defer logfile.Close()
	logrus.New()
	logrus.SetOutput(logfile)
	logrus.SetFormatter(new(logrus.JSONFormatter))
	m, _ := repo.NewMiner(&repo.Source{Resources: viper.GetString("agent.metricsType")})
	serv := service.New(m)
	newAgent := agent.NewAgent(serv)

	defer newAgent.Stop()
	if err := newAgent.RunWithSerialized(); err != nil {
		logrus.Errorf("error while running agent: %s", err.Error())
	}
}
