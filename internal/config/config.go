package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Database Database
	JWT      JWT
	Logger   Logger
}

type Server struct {
	Url string `mapstructure:"url"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWT struct {
	SecretKey      string `mapstructure:"secret_key"`
	ExpirationTime int    `mapstructure:"expiration_time"`
}

type Logger struct {
	Sink  string `mapstructure:"sink"`
	Level string `mapstructure:"level"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("AVITO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading configuration: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing configuration: %w", err)
	}

	return &cfg, nil
}
