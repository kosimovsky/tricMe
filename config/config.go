package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitServerConfig reads configuration file and ENV for server
func InitServerConfig() error {
	viper.New()
	viper.AddConfigPath("config")
	viper.SetConfigName(".server")
	viper.SetConfigType("yaml")

	viper.SetDefault("Address", "127.0.0.1:8080")
	viper.SetDefault("Restore", true)
	viper.SetDefault("Interval", 300)
	viper.SetDefault("File", "/tmp/devops-metrics-db.json")
	viper.SetDefault("Logfile", "server.log")
	viper.SetDefault("Loglevel", 3)
	viper.SetDefault("GinMode", "release")
	viper.SetDefault("Debug", false)
	viper.SetDefault("Storage", "memory")

	fSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	fSet.StringP("address", "a", "127.0.0.1:8080", "address for server")
	fSet.BoolP("restore", "r", true, "restore metrics from file")
	fSet.IntP("interval", "i", 300, "interval for storing metrics to file")
	fSet.StringP("file", "f", "/tmp/devops-metrics-db.json", "file to store metrics")
	err := viper.BindPFlags(fSet)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if err = fSet.Parse(os.Args[1:]); err != nil {
		logrus.Error(err.Error())
		return err
	}

	if addr := os.Getenv("ADDRESS"); addr != "" {
		viper.Set("Address", addr)
	}
	if restore := os.Getenv("RESTORE"); restore != "" {
		viper.Set("Restore", restore)
	}
	if interval := os.Getenv("STORE_INTERVAL"); interval != "" {
		viper.Set("Interval", interval)
	}
	if file := os.Getenv("STORE_FILE"); file != "" {
		viper.Set("File", file)
	}

	if err := viper.ReadInConfig(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Use config file:", viper.ConfigFileUsed())
		if err != nil {
			return err
		}
	}

	setGinMode(viper.GetString("GinMode"))
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

// setGinMode sets mode for Gin Framework (default: debug)
func setGinMode(mode string) {
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// InitAgentConfig reads configuration file and ENV for agent
func InitAgentConfig() error {
	viper.New()
	viper.AddConfigPath("config")
	viper.SetConfigName(".agent")
	viper.SetConfigType("yaml")

	viper.SetDefault("Address", "127.0.0.1:8080")
	viper.SetDefault("Poll", "10s")
	viper.SetDefault("Report", "2s")
	viper.SetDefault("Logfile", "agent.log")
	viper.SetDefault("MetricsType", "memStat")

	fSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	fSet.StringP("address", "a", "127.0.0.1:8080",
		`server address to which metrics should be sent`)
	fSet.StringP("poll", "r", "10s", "report interval")
	fSet.StringP("report", "p", "2s", "poll interval")

	err := viper.BindPFlags(fSet)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if err = fSet.Parse(os.Args[1:]); err != nil {
		logrus.Error(err.Error())
		return err
	}

	if addr := os.Getenv("ADDRESS"); addr != "" {
		viper.Set("Address", addr)
	}
	if poll := os.Getenv("POLL_INTERVAL"); poll != "" {
		viper.Set("Poll", poll)
	}
	if report := os.Getenv("REPORT_INTERVAL"); report != "" {
		viper.Set("Report", report)
	}

	if err = viper.ReadInConfig(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Use config file:", viper.ConfigFileUsed())
		if err != nil {
			return err
		}
	}
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
