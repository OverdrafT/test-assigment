package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"SERVER_PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string
	DB_NAME     string
	DB_HOST     string
}

func LoadConfig(name string) (*Config, error) {
	viper.AutomaticEnv()
	viper.AddConfigPath("./")
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read .env file: ", err)

	}

	var cfg Config
	err = viper.Unmarshal(&cfg)

	return &cfg, err
}
