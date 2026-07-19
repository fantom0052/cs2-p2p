package market

import (
	"strings"
	"time"
)

// MarketName идентифицирует торговую площадку внутри нашей системы.
type MarketName string

// MarketHashName хранит каноничное имя предмета, по которому мы связываем данные разных маркетов.
type MarketHashName string

// ItemType описывает категорию предмета на торговой площадке: оружие, кейс, стикер, агент и т.д.
type ItemType string

// ImageURL хранит ссылку на изображение предмета, полученную с торговой площадки.
type ImageURL string

// Currency хранит ISO 4217 код валюты, например USD, EUR или RUB.
type Currency string

// Money хранит денежное значение в минимальных единицах валюты, чтобы избежать ошибок округления float.
type Money struct {
	AmountMinor int64
	Currency    Currency
}

// MarketItemSnapshot описывает состояние предмета на маркете в конкретный момент времени.
type MarketItemSnapshot struct {
	Market         MarketName
	MarketHashName MarketHashName
	Price          Money
	AvailableCount int
	ItemType       ItemType
	ImageURL       ImageURL
	CollectedAt    time.Time
}

const MarketSteam MarketName = "steam"

// NewMarketName нормализует и валидирует название торговой площадки.
func NewMarketName(value string) (MarketName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyMarketName
	}

	return MarketName(strings.ToLower(value)), nil
}

// NewMarketHashName валидирует каноничное market hash name предмета.
func NewMarketHashName(value string) (MarketHashName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyMarketHashName
	}

	return MarketHashName(value), nil
}

// NewItemType валидирует категорию предмета на торговой площадке.
func NewItemType(value string) (ItemType, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyItemType
	}

	return ItemType(value), nil
}

// NewImageURL валидирует ссылку на изображение предмета.
func NewImageURL(value string) (ImageURL, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyImageURL
	}

	return ImageURL(value), nil
}

// NewCurrency нормализует и валидирует ISO 4217 код валюты.
func NewCurrency(value string) (Currency, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if len(value) != 3 {
		return "", ErrInvalidCurrency
	}

	return Currency(value), nil
}

// NewMoney валидирует и создает денежное значение в минимальных единицах валюты.
func NewMoney(amountMinor int64, currency Currency) (Money, error) {
	if amountMinor < 0 {
		return Money{}, ErrInvalidMoneyAmount
	}
	if len(currency) != 3 {
		return Money{}, ErrInvalidCurrency
	}

	return Money{AmountMinor: amountMinor, Currency: currency}, nil
}

// NewMarketItemSnapshot валидирует и создает наблюдение за состоянием предмета на маркете.
func NewMarketItemSnapshot(
	marketName MarketName,
	marketHashName MarketHashName,
	price Money,
	availableCount int,
	itemType ItemType,
	imageURL ImageURL,
	collectedAt time.Time,
) (MarketItemSnapshot, error) {
	if marketName == "" {
		return MarketItemSnapshot{}, ErrEmptyMarketName
	}
	if marketHashName == "" {
		return MarketItemSnapshot{}, ErrEmptyMarketHashName
	}
	if len(price.Currency) != 3 {
		return MarketItemSnapshot{}, ErrInvalidCurrency
	}
	if price.AmountMinor <= 0 {
		return MarketItemSnapshot{}, ErrInvalidSnapshotPrice
	}
	if availableCount < 0 {
		return MarketItemSnapshot{}, ErrNegativeListingCount
	}
	if itemType == "" {
		return MarketItemSnapshot{}, ErrEmptyItemType
	}
	if imageURL == "" {
		return MarketItemSnapshot{}, ErrEmptyImageURL
	}
	if collectedAt.IsZero() {
		return MarketItemSnapshot{}, ErrEmptySnapshotTime
	}

	return MarketItemSnapshot{
		Market:         marketName,
		MarketHashName: marketHashName,
		Price:          price,
		AvailableCount: availableCount,
		ItemType:       itemType,
		ImageURL:       imageURL,
		CollectedAt:    collectedAt.UTC(),
	}, nil
}
