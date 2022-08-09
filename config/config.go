package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)

// InitServerConfig reads configuration file and ENV for server
func InitServerConfig() error {
	viper.New()
	viper.AddConfigPath("config")
	viper.SetConfigName(".server")
	viper.SetConfigType("yaml")

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

	fSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	addr := fSet.StringP("address", "a", "127.0.0.1:8080", "address for server")
	restore := fSet.BoolP("restore", "r", true, "restore metrics from file")
	interval := fSet.DurationP("interval", "i", 300, "interval for storing metrics to file")
	file := fSet.StringP("file", "f", "/tmp/devops-metrics-db.json", "file to store metrics")
	err := viper.BindPFlags(fSet)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if err = fSet.Parse(os.Args[1:]); err != nil {
		logrus.Error(err.Error())
		return err
	}

	if viper.GetString("Address") == "" {
		viper.Set("Address", addr)
	}
	if viper.GetString("Restore") == "" {
		viper.Set("Address", restore)
	}
	if viper.GetString("Interval") == "" {
		viper.Set("Interval", interval)
	}
	if viper.GetString("File") == "" {
		viper.Set("File", file)
	}

	if err := viper.ReadInConfig(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Use serverConfig file:", viper.ConfigFileUsed())
		if err != nil {
			return err
		}
	}

	setGinMode(viper.GetString("GinMode"))
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

	if addr := os.Getenv("ADDRESS"); addr != "" {
		viper.Set("Address", addr)
	}
	if poll := os.Getenv("POLL_INTERVAL"); poll != "" {
		viper.Set("Poll", poll)
	}
	if report := os.Getenv("REPORT_INTERVAL"); report != "" {
		viper.Set("Report", report)
	}

	fSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	addr := fSet.StringP("address", "a", "127.0.0.1:8080",
		`server address to which metrics should be sent`)
	poll := fSet.DurationP("poll", "r", 2*time.Second, "report interval")
	report := fSet.DurationP("report", "p", 10*time.Second, "poll interval")

	err := viper.BindPFlags(fSet)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if err = fSet.Parse(os.Args[1:]); err != nil {
		logrus.Error(err.Error())
		return err
	}

	if viper.GetString("Address") == "" {
		viper.Set("Address", *addr)
	}
	if viper.GetString("Poll") == "" {
		viper.Set("Poll", *poll)
	}
	if viper.GetString("Report") == "" {
		viper.Set("Report", *report)
	}

	if err = viper.ReadInConfig(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Use config file:", viper.ConfigFileUsed())
		if err != nil {
			return err
		}
	}
	return nil
}

type serverConfig struct {
	Address       string
	Debug         bool
	GinMode       string
	Logfile       string
	Loglevel      int
	Storage       string
	StoreInterval time.Duration
	Filename      string
	Restore       bool
}

type agentConfig struct {
	Address     string
	Logfile     string
	MetricsType string
	Poll        time.Duration
	Report      time.Duration
}

func AgentConfig() *agentConfig {
	return &agentConfig{
		Address:     viper.GetString("Address"),
		Logfile:     viper.GetString("Logfile"),
		MetricsType: viper.GetString("MetricsType"),
		Poll:        viper.GetDuration("Poll"),
		Report:      viper.GetDuration("Report"),
	}
}

func ServerConfig() *serverConfig {
	return &serverConfig{
		Address:       viper.GetString("Address"),
		Debug:         viper.GetBool("Debug"),
		GinMode:       viper.GetString("ginMode"),
		Logfile:       viper.GetString("Logfile"),
		Loglevel:      viper.GetInt("Loglevel"),
		Storage:       viper.GetString("Storage"),
		StoreInterval: viper.GetDuration("Interval"),
		Filename:      viper.GetString("File"),
		Restore:       viper.GetBool("Restore"),
	}
}
