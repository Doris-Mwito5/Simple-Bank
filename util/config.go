package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER      string `mapstructure:"DB_DRIVER"`
	DB_SOURCE      string `mapstructure:"DB_SOURCE"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	
	// Set up viper to look for the configuration file
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Read environment variables with AutomaticEnv
	viper.AutomaticEnv()

	// Read in the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	// Unmarshal the config into the Config struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
