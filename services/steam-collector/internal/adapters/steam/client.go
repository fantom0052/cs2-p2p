package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fanto/p2p/steam-collector/internal/domain/market"
	"github.com/fanto/p2p/steam-collector/internal/usecase/collector"
)

// Client получает данные о предметах с торговой площадки Steam Market.
type Client struct {
	httpClient *http.Client
	options    ClientOptions

	cookieMu           sync.Mutex
	cookiesInitialized bool
}

// ClientOptions содержит настройки, необходимые Steam-клиенту для выполнения запросов.
type ClientOptions struct {
	HomeURL            string
	BaseURL            string
	ImageBaseURL       string
	AppID              int
	CurrencyID         int
	Currency           market.Currency
	UserAgent          string
	AcceptLanguage     string
	CookieSettings     string
	RequestTimeout     time.Duration
	NoRender           int
	DefaultSortColumn  string
	DefaultSortDir     string
	DefaultPageSize    int
	SearchDescriptions int
	BootstrapCookies   bool
	CategoryFilters    map[string]string
}

// NewClient создает Steam adapter.
func NewClient(httpClient *http.Client, options ClientOptions) (*Client, error) {
	if err := validateClientOptions(options); err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: options.RequestTimeout}
	}
	if httpClient.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("create steam cookie jar: %w", err)
		}
		httpClient.Jar = jar
	}

	return &Client{
		httpClient: httpClient,
		options:    options,
	}, nil
}

// FetchMarketItems получает одну страницу предметов Steam Market и возвращает ее в domain-модели.
func (c *Client) FetchMarketItems(ctx context.Context, params collector.FetchParams) (collector.FetchPage, error) {
	if err := c.ensureCookies(ctx); err != nil {
		return collector.FetchPage{}, err
	}

	requestURL, err := c.buildSearchURL(params)
	if err != nil {
		return collector.FetchPage{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return collector.FetchPage{}, err
	}
	c.setBrowserHeaders(req, "empty", "cors", "same-origin")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", c.options.HomeURL+"market/")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return collector.FetchPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return collector.FetchPage{}, collector.ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return collector.FetchPage{}, fmt.Errorf("steam search returned status %d", resp.StatusCode)
	}

	var dto searchRenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return collector.FetchPage{}, err
	}
	if !dto.Success {
		return collector.FetchPage{}, ErrSearchResponseFailed
	}

	collectedAt := time.Now().UTC()
	items := make([]market.MarketItemSnapshot, 0, len(dto.Results))
	for _, item := range dto.Results {
		snapshot, err := mapSearchRenderItem(item, c.options.Currency, c.options.ImageBaseURL, collectedAt)
		if err != nil {
			return collector.FetchPage{}, err
		}
		items = append(items, snapshot)
	}

	return collector.FetchPage{
		Items:      items,
		Start:      dto.Start,
		PageSize:   dto.PageSize,
		TotalCount: dto.TotalCount,
	}, nil
}

func (c *Client) ensureCookies(ctx context.Context) error {
	if !c.options.BootstrapCookies {
		return nil
	}

	c.cookieMu.Lock()
	defer c.cookieMu.Unlock()

	if c.cookiesInitialized {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.options.HomeURL, nil)
	if err != nil {
		return err
	}
	c.setBrowserHeaders(req, "document", "navigate", "none")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return collector.ErrRateLimited
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("steam bootstrap returned status %d", resp.StatusCode)
	}

	c.addStaticCookies()
	c.cookiesInitialized = true
	return nil
}

func (c *Client) addStaticCookies() {
	parsed, err := url.Parse(c.options.HomeURL)
	if err != nil {
		return
	}

	cookies := []*http.Cookie{
		{Name: "Steam_Language", Value: "english"},
		{Name: "timezoneOffset", Value: "10800,0"},
	}
	if c.options.CookieSettings != "" {
		cookies = append(cookies, &http.Cookie{Name: "cookieSettings", Value: c.options.CookieSettings})
	}

	c.httpClient.Jar.SetCookies(parsed, cookies)
}

func (c *Client) setBrowserHeaders(req *http.Request, fetchDest, fetchMode, fetchSite string) {
	req.Header.Set("Accept-Language", c.options.AcceptLanguage)
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", fetchDest)
	req.Header.Set("Sec-Fetch-Mode", fetchMode)
	req.Header.Set("Sec-Fetch-Site", fetchSite)
	req.Header.Set("User-Agent", c.options.UserAgent)
}

func (c *Client) buildSearchURL(params collector.FetchParams) (string, error) {
	params = c.normalizeSearchParams(params)

	parsed, err := url.Parse(c.options.BaseURL)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Set("norender", strconv.Itoa(c.options.NoRender))
	query.Set("appid", strconv.Itoa(c.options.AppID))
	query.Set("currency", strconv.Itoa(c.options.CurrencyID))
	query.Set("query", params.Query)
	query.Set("search_descriptions", strconv.Itoa(c.options.SearchDescriptions))
	query.Set("sort_column", params.SortColumn)
	query.Set("sort_dir", params.SortDir)
	query.Set("start", strconv.Itoa(params.Start))
	query.Set("count", strconv.Itoa(params.Count))

	for key, value := range c.options.CategoryFilters {
		if value == "" || value == "any" {
			continue
		}
		query.Set(key, value)
	}

	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (c *Client) normalizeSearchParams(params collector.FetchParams) collector.FetchParams {
	if params.SortColumn == "" {
		params.SortColumn = c.options.DefaultSortColumn
	}
	if params.SortDir == "" {
		params.SortDir = c.options.DefaultSortDir
	}
	if params.Count <= 0 {
		params.Count = c.options.DefaultPageSize
	}
	if params.Start < 0 {
		params.Start = 0
	}

	return params
}

func validateClientOptions(options ClientOptions) error {
	if options.HomeURL == "" {
		return ErrEmptyHomeURL
	}
	if !strings.HasSuffix(options.HomeURL, "/") {
		return ErrInvalidHomeURL
	}
	if options.BaseURL == "" {
		return ErrEmptyBaseURL
	}
	if options.ImageBaseURL == "" {
		return ErrEmptyImageBaseURL
	}
	if options.AppID <= 0 {
		return ErrInvalidAppID
	}
	if options.CurrencyID <= 0 {
		return ErrInvalidCurrencyID
	}
	if len(options.Currency) != 3 {
		return market.ErrInvalidCurrency
	}
	if options.UserAgent == "" {
		return ErrEmptyUserAgent
	}
	if options.AcceptLanguage == "" {
		return ErrEmptyAcceptLanguage
	}
	if options.RequestTimeout <= 0 {
		return ErrInvalidRequestTimeout
	}
	if options.NoRender != 0 && options.NoRender != 1 {
		return ErrInvalidNoRender
	}
	if options.DefaultSortColumn == "" {
		return ErrEmptyDefaultSortColumn
	}
	if options.DefaultSortDir == "" {
		return ErrEmptyDefaultSortDir
	}
	if options.DefaultPageSize <= 0 {
		return ErrInvalidDefaultPageSize
	}
	if options.SearchDescriptions != 0 && options.SearchDescriptions != 1 {
		return ErrInvalidSearchDescriptions
	}

	return nil
}
