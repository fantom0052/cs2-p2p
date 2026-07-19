package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит настройки приложения.
type Config struct {
	LisSkins LisSkins
	Kafka    Kafka
}

// LisSkins содержит настройки для dump-а и WebSocket Lis-Skins.
type LisSkins struct {
	APIKey         string        `envconfig:"LIS_SKINS_API_KEY"`
	DumpURL        string        `envconfig:"LIS_SKINS_DUMP_URL" default:"https://lis-skins.com/market_export_json/api_csgo_full.json"`
	WebSocketURL   string        `envconfig:"LIS_SKINS_WS_CON_URL" default:"wss://ws.lis-skins.com/connection/websocket"`
	WSTokenURL     string        `envconfig:"LIS_SKINS_WS_TOKEN_URL" default:"https://api.lis-skins.com/v1/user/get-ws-token"`
	Currency       string        `envconfig:"LIS_SKINS_CURRENCY" default:"USD"`
	RequestTimeout time.Duration `envconfig:"LIS_SKINS_REQUEST_TIMEOUT" default:"5m"`
	BatchSize      int           `envconfig:"LIS_SKINS_DUMP_BATCH_SIZE" default:"2000"`
	DumpLimit      int           `envconfig:"LIS_SKINS_DUMP_LIMIT" default:"10000"`
}

// Kafka содержит настройки для публикации сообщений в Kafka-compatible broker.
type Kafka struct {
	Brokers       []string      `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	Topic         string        `envconfig:"KAFKA_TOPIC" default:"raw.market.prices.lisskins"`
	ClientID      string        `envconfig:"KAFKA_CLIENT_ID" default:"lisskins-collector"`
	Acks          string        `envconfig:"KAFKA_ACKS" default:"all"`
	Compression   string        `envconfig:"KAFKA_COMPRESSION" default:"snappy"`
	Linger        time.Duration `envconfig:"KAFKA_LINGER" default:"100ms"`
	BatchBytes    int32         `envconfig:"KAFKA_BATCH_BYTES" default:"1048576"`
	RecordRetries int           `envconfig:"KAFKA_RECORD_RETRIES" default:"3"`
}

// NewConfig читает настройки приложения из переменных окружения.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("load config from env: %w", err)
	}

	return cfg, nil
}
