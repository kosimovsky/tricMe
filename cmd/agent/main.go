package main

import (
	"fmt"
	"github.com/kosimovsky/tricMe/internal/agent"
	"github.com/kosimovsky/tricMe/internal/repo"
	"github.com/kosimovsky/tricMe/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		_ = fmt.Errorf("error while reading config file %v", err.Error())
	}

	logfileFromConfig := viper.GetString("agent.logfile")
	logfile, err := os.OpenFile(logfileFromConfig, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error occured opening file : %v, %s", err, logfileFromConfig)
	}
	defer logfile.Close()
	logrus.New()
	logrus.SetOutput(logfile)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	m, _ := repo.NewMiner(&repo.Source{Resources: "memStat"})
	serv := service.New(m)
	newAgent := agent.NewAgent(serv)
	logrus.Fatalln(newAgent.Run())
}

// initConfig reads configuration file
func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".config")

	if err := viper.ReadInConfig(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Use config file:", viper.ConfigFileUsed())
		if err != nil {
			return err
		}
	}
	return nil
}
