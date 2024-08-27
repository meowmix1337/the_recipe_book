package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Hostname    string `mapstructure:"HOSTNAME"`
	Port        string `mapstructure:"PORT"`
}

// NewConfig initializes and returns a Config instance
func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("env") // Use .env files for environment configuration

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("HOSTNAME", "localhost")
	viper.SetDefault("PORT", "8080")

	err := viper.ReadInConfig() // Read from config file
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("Error reading config file: %v. Using defaults and environment variables.", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
