package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config interface {
	GetEnvironment() string
	GetJWTSecret() string
	GetPort() string
}

// Config holds the application configuration.
type ConfigImpl struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Hostname    string `mapstructure:"HOSTNAME"`
	Port        string `mapstructure:"PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

var _ Config = (*ConfigImpl)(nil)

// NewConfig initializes and returns a Config instance.
func NewConfig() (*ConfigImpl, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("env") // Use .env files for environment configuration.

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("HOSTNAME", "localhost")
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("LOG_LEVEL", "debug")
	// You should definitely replace with your own secret, this is for testing only
	viper.SetDefault("JWT_SECRET", "ZY7FM2SuRM13eRYX")

	err := viper.ReadInConfig() // Read from config file.
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("Error reading config file: %v. Using defaults and environment variables.", err))
	}

	var config ConfigImpl
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *ConfigImpl) GetEnvironment() string {
	return c.Environment
}

func (c *ConfigImpl) GetJWTSecret() string {
	return c.JWTSecret
}

func (c *ConfigImpl) GetPort() string {
	return c.Port
}
