package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит настройки приложения.
type Config struct {
	Steam     Steam
	Kafka     Kafka
	Scheduler Scheduler
}

// Steam содержит настройки для работы со Steam Market в простых Go-типах.
type Steam struct {
	HomeURL                string        `envconfig:"STEAM_HOME_URL" default:"https://steamcommunity.com/"`
	BaseURL                string        `envconfig:"STEAM_BASE_URL" default:"https://steamcommunity.com/market/search/render/"`
	ImageBaseURL           string        `envconfig:"STEAM_IMAGE_BASE_URL" default:"https://community.cloudflare.steamstatic.com/economy/image/"`
	AppID                  int           `envconfig:"STEAM_APP_ID" default:"730"`
	CurrencyID             int           `envconfig:"STEAM_CURRENCY_ID" default:"1"`
	Currency               string        `envconfig:"STEAM_CURRENCY" default:"USD"`
	UserAgent              string        `envconfig:"STEAM_USER_AGENT" default:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36"`
	AcceptLanguage         string        `envconfig:"STEAM_ACCEPT_LANGUAGE" default:"ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"`
	CookieSettings         string        `envconfig:"STEAM_COOKIE_SETTINGS" default:"%7B%22version%22%3A1%2C%22preference_state%22%3A1%2C%22content_customization%22%3Anull%2C%22valve_analytics%22%3Anull%2C%22third_party_analytics%22%3Anull%2C%22third_party_content%22%3Anull%2C%22utm_enabled%22%3Atrue%7D"`
	BootstrapCookies       bool          `envconfig:"STEAM_BOOTSTRAP_COOKIES" default:"true"`
	RequestTimeout         time.Duration `envconfig:"STEAM_REQUEST_TIMEOUT" default:"15s"`
	NoRender               int           `envconfig:"STEAM_NORENDER" default:"1"`
	DefaultSortColumn      string        `envconfig:"STEAM_DEFAULT_SORT_COLUMN" default:"popular"`
	DefaultSortDir         string        `envconfig:"STEAM_DEFAULT_SORT_DIR" default:"desc"`
	DefaultPageSize        int           `envconfig:"STEAM_DEFAULT_PAGE_SIZE" default:"30"`
	SearchDescriptions     int           `envconfig:"STEAM_SEARCH_DESCRIPTIONS" default:"0"`
	CategoryItemSet        string        `envconfig:"STEAM_CATEGORY_ITEM_SET" default:"any"`
	CategoryProPlayer      string        `envconfig:"STEAM_CATEGORY_PRO_PLAYER" default:"any"`
	CategoryStickerCapsule string        `envconfig:"STEAM_CATEGORY_STICKER_CAPSULE" default:"any"`
	CategoryTournament     string        `envconfig:"STEAM_CATEGORY_TOURNAMENT" default:"any"`
	CategoryTournamentTeam string        `envconfig:"STEAM_CATEGORY_TOURNAMENT_TEAM" default:"any"`
	CategoryType           string        `envconfig:"STEAM_CATEGORY_TYPE" default:"any"`
	CategoryWeapon         string        `envconfig:"STEAM_CATEGORY_WEAPON" default:"any"`
}

// Kafka содержит настройки для публикации сообщений в Kafka-compatible broker.
type Kafka struct {
	Brokers       []string      `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	Topic         string        `envconfig:"KAFKA_TOPIC" default:"raw.market.prices.steam"`
	ClientID      string        `envconfig:"KAFKA_CLIENT_ID" default:"steam-collector"`
	Acks          string        `envconfig:"KAFKA_ACKS" default:"all"`
	Compression   string        `envconfig:"KAFKA_COMPRESSION" default:"snappy"`
	Linger        time.Duration `envconfig:"KAFKA_LINGER" default:"100ms"`
	BatchBytes    int32         `envconfig:"KAFKA_BATCH_BYTES" default:"1048576"`
	RecordRetries int           `envconfig:"KAFKA_RECORD_RETRIES" default:"3"`
}

// Scheduler содержит настройки цикла сбора данных.
type Scheduler struct {
	// Interval задает паузу между запросами страниц Steam Market.
	Interval time.Duration `envconfig:"SCHEDULER_INTERVAL" default:"10s"`
	// RateLimitPause задает паузу после ограничения частоты запросов источником.
	RateLimitPause time.Duration `envconfig:"SCHEDULER_RATE_LIMIT_PAUSE" default:"10m"`
	// Jitter добавляет случайную задержку к каждой паузе между запросами.
	Jitter time.Duration `envconfig:"SCHEDULER_JITTER" default:"10s"`
	// PageSize задает размер страницы, которую scheduler просит собрать за одну итерацию.
	PageSize int `envconfig:"SCHEDULER_PAGE_SIZE" default:"30"`
	// Query ограничивает сбор предметами, подходящими под поисковую строку Steam Market.
	Query string `envconfig:"SCHEDULER_QUERY" default:""`
}

// NewConfig читает настройки приложения из переменных окружения.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("load config from env: %w", err)
	}

	return cfg, nil
}
