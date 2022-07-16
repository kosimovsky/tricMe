package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig reads configuration file
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
	return nil
}
