package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitConfig reads configuration file and ENV
func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".config")

	if err := viper.ReadInConfig(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Use config file:", viper.ConfigFileUsed())
		if err != nil {
			return err
		}
	}

	err := viper.BindEnv("server.address", "ADDRESS")
	if err != nil {
		return err
	}
	err = viper.BindEnv("agent.pollInterval", "POLL_INTERVAL")
	if err != nil {
		return err
	}
	err = viper.BindEnv("agent.reportInterval", "REPORT_INTERVAL")
	if err != nil {
		return err
	}
	err = viper.BindEnv("server.store.storeInterval", "STORE_INTERVAL")
	if err != nil {
		return err
	}
	err = viper.BindEnv("server.store.storeFile", "STORE_FILE")
	if err != nil {
		return err
	}
	err = viper.BindEnv("server.store.restore", "RESTORE")
	if err != nil {
		return err
	}

	setGinMode(viper.GetString("server.ginMode"))
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
