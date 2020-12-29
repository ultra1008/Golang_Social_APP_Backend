package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	Login    string `envconfig:"DB_LOGIN"`
	Password string `envconfig:"DB_PASSWORD"`
	DBName   string `envconfig:"DB_NAME"`
}

func (d *DBConfig) ConnectionString() string {
	return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", d.Login, d.Password, d.Host, d.Port, d.DBName)
}

func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("creating config: %w", err)
	}

	return &cfg, nil
}
