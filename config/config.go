package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var SessionName = "hsn-session"

type Config struct {
	DB        *DBConfig
	Server    *HTTPServerConfig
	SecretKey string `envconfig:"SESSION_SECRET_KEY" default:"verysecretkey"`
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"3306"`
	Login    string `envconfig:"DB_LOGIN" default:"niklod"`
	Password string `envconfig:"DB_PASSWORD" default:"VLQi4Vttuo6wFRqm"`
	DBName   string `envconfig:"DB_NAME" default:"hsn"`
}

func (d *DBConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", d.Login, d.Password, d.Host, d.Port, d.DBName)
}

type RabbitMQConfig struct {
	Host             string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	Port             string `envconfig:"RABBITMQ_PORT" default:"5672"`
	Login            string `envconfig:"RABBITMQ_USERNAME" default:""`
	Password         string `envconfig:"RABBITMQ_PASSWORD" default:""`
	FeedQueueName    string `envconfig:"RABBITMQ_FEED_QUEUE_NAME" default:"feedQueue"`
	FeedExchangeName string `envconfig:"RABBITMQ_FEED_EXCHANGE_NAME" default:"feedExchange"`
	FeedRoutingKey   string `envconfig:"RABBITMQ_FEED_ROUTING_KEY" default:"feedUpdate"`
	ReceiversCount   int    `envconfig:"RABBITMQ_FEED_RECEIVERS_COUNT" default:"2"`
}

func (r *RabbitMQConfig) ConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.Login, r.Password, r.Host, r.Port)
}

type HTTPServerConfig struct {
	Port int `envconfig:"HTTP_SERVER_PORT" default:"8080"`
}

func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("creating config: %w", err)
	}

	return &cfg, nil
}
