package lisskins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/centrifugal/centrifuge-go"

	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
	"github.com/fanto/p2p/lisskins-collector/internal/usecase/collector"
)

// Client получает начальный dump и live-события Lis-Skins.
type Client struct {
	httpClient *http.Client
	options    ClientOptions
	token      *tokenManager
}

// ClientOptions содержит настройки Lis-Skins adapter-а.
type ClientOptions struct {
	APIKey       string
	DumpURL      string
	WebSocketURL string
	WSTokenURL   string
	Currency     market.Currency
	DumpLimit    int
}

// NewClient создает adapter для Lis-Skins.
func NewClient(httpClient *http.Client, options ClientOptions) (*Client, error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: httpClient,
		options:    options,
		token:      newTokenManager(httpClient, options.APIKey, options.WSTokenURL, 2*time.Minute),
	}, nil
}

// LoadDump загружает полный dump Lis-Skins и приводит его к domain snapshot-ам.
func (c *Client) LoadDump(ctx context.Context) ([]market.MarketItemSnapshot, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.options.DumpURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("build dump request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send dump request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lisskins dump returned status %d", resp.StatusCode)
	}

	var body dumpResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode dump response: %w", err)
	}

	collectedAt := time.Now().UTC()
	snapshots := make([]market.MarketItemSnapshot, 0, len(body.Items))
	for i, item := range body.Items {
		if c.options.DumpLimit > 0 && i >= c.options.DumpLimit {
			break
		}

		snapshot, err := mapDumpItem(item, c.options.Currency, collectedAt)
		if err != nil {
			continue
		}
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// Watch подключается к Lis-Skins WebSocket и передает live snapshot-ы в handler.
func (c *Client) Watch(ctx context.Context, handler collector.SnapshotHandler) error {
	wsClient := centrifuge.NewJsonClient(c.options.WebSocketURL, centrifuge.Config{
		GetToken: func(_ centrifuge.ConnectionTokenEvent) (string, error) {
			return c.token.token(ctx)
		},
	})

	sub, err := wsClient.NewSubscription("public:obtained-skins")
	if err != nil {
		return fmt.Errorf("create lisskins subscription: %w", err)
	}

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		var event skinEvent
		if err := json.Unmarshal(e.Data, &event); err != nil {
			return
		}

		snapshot, err := mapSkinEvent(event, c.options.Currency, time.Now().UTC())
		if err != nil {
			return
		}

		_ = handler(ctx, snapshot)
	})

	if err := sub.Subscribe(); err != nil {
		return fmt.Errorf("subscribe lisskins websocket: %w", err)
	}
	if err := wsClient.Connect(); err != nil {
		return fmt.Errorf("connect lisskins websocket: %w", err)
	}

	<-ctx.Done()
	_ = wsClient.Disconnect()
	wsClient.Close()
	return ctx.Err()
}

func validateOptions(options ClientOptions) error {
	if options.DumpURL == "" {
		return ErrEmptyDumpURL
	}
	if options.WebSocketURL == "" {
		return ErrEmptyWebSocketURL
	}
	if options.WSTokenURL == "" {
		return ErrEmptyTokenURL
	}
	if options.APIKey == "" {
		return ErrEmptyAPIKey
	}
	if len(options.Currency) != 3 {
		return market.ErrInvalidCurrency
	}

	return nil
}

type tokenManager struct {
	httpClient *http.Client
	apiKey     string
	tokenURL   string
	ttl        time.Duration
	mu         sync.RWMutex
	cached     string
	expiresAt  time.Time
}

func newTokenManager(httpClient *http.Client, apiKey string, tokenURL string, ttl time.Duration) *tokenManager {
	return &tokenManager{
		httpClient: httpClient,
		apiKey:     apiKey,
		tokenURL:   tokenURL,
		ttl:        ttl,
	}
}

func (tm *tokenManager) token(ctx context.Context) (string, error) {
	now := time.Now()

	tm.mu.RLock()
	if tm.cached != "" && now.Before(tm.expiresAt) {
		token := tm.cached
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	token, err := tm.fetch(ctx)
	if err != nil {
		return "", err
	}

	tm.mu.Lock()
	tm.cached = token
	tm.expiresAt = now.Add(tm.ttl)
	tm.mu.Unlock()

	return token, nil
}

func (tm *tokenManager) fetch(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tm.tokenURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("build websocket token request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tm.apiKey)

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send websocket token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("websocket token request returned status %d", resp.StatusCode)
	}

	var body tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode websocket token response: %w", err)
	}

	return body.Data.Token, nil
}
