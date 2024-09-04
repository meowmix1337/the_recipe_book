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

	GetDBUser() string
	GetDBPassword() string
	GetDBName() string
	GetDBHost() string
	GetDBPort() string
}

// Config holds the application configuration.
type ConfigImpl struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Hostname    string `mapstructure:"HOSTNAME"`
	Port        string `mapstructure:"PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_POST"`
	DBName      string `mapstructure:"DB_NAME"`
}

var _ Config = (*ConfigImpl)(nil)

// NewConfig initializes and returns a Config instance.
func NewConfig() (*ConfigImpl, error) {
	viper.SetConfigName(".env")
	viper.AddConfigPath("../") // Specify the root directory for the config
	viper.SetConfigType("env") // Use .env files for environment configuration.

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("HOSTNAME", "localhost")
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("LOG_LEVEL", "debug")
	// You should definitely replace with your own secret, this is for testing only
	viper.SetDefault("JWT_SECRET", "some_really_bad_secret")

	// Database
	viper.SetDefault("DB_USER", "admin")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "the_recipe_book")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")

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

func (c *ConfigImpl) GetDBUser() string {
	return c.DBUser
}

func (c *ConfigImpl) GetDBPassword() string {
	return c.DBPassword
}

func (c *ConfigImpl) GetDBName() string {
	return c.DBName
}

func (c *ConfigImpl) GetDBHost() string {
	return c.DBHost
}

func (c *ConfigImpl) GetDBPort() string {
	return c.DBPort
}
