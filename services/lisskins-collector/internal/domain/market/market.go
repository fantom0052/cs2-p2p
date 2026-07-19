package market

import (
	"strings"
	"time"
)

// MarketName идентифицирует торговую площадку внутри нашей системы.
type MarketName string

// ExternalItemID хранит ID конкретного листинга на стороне площадки.
type ExternalItemID string

// MarketHashName хранит каноничное имя предмета, по которому мы связываем данные разных маркетов.
type MarketHashName string

// ItemType описывает категорию предмета на торговой площадке.
type ItemType string

// Currency хранит ISO 4217 код валюты, например USD, EUR или RUB.
type Currency string

// SnapshotEvent показывает, какое состояние предмета пришло от источника.
type SnapshotEvent string

// Money хранит денежное значение в минимальных единицах валюты, чтобы избежать ошибок округления float.
type Money struct {
	AmountMinor int64
	Currency    Currency
}

// MarketItemSnapshot описывает состояние конкретного листинга на маркете в конкретный момент времени.
type MarketItemSnapshot struct {
	Market         MarketName
	ExternalItemID ExternalItemID
	MarketHashName MarketHashName
	Price          Money
	AvailableCount int
	ItemType       ItemType
	Event          SnapshotEvent
	CollectedAt    time.Time
}

const (
	MarketLisSkins MarketName = "lisskins"

	EventSnapshot     SnapshotEvent = "snapshot"
	EventItemAdded    SnapshotEvent = "item_added"
	EventPriceChanged SnapshotEvent = "price_changed"
	EventItemRemoved  SnapshotEvent = "item_removed"
)

// NewExternalItemID валидирует ID листинга из внешней площадки.
func NewExternalItemID(value string) (ExternalItemID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyExternalItemID
	}

	return ExternalItemID(value), nil
}

// NewMarketHashName валидирует каноничное market hash name предмета.
func NewMarketHashName(value string) (MarketHashName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyMarketHashName
	}

	return MarketHashName(value), nil
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

// NewMarketItemSnapshot валидирует и создает snapshot конкретного листинга.
func NewMarketItemSnapshot(
	marketName MarketName,
	externalItemID ExternalItemID,
	marketHashName MarketHashName,
	price Money,
	availableCount int,
	itemType ItemType,
	event SnapshotEvent,
	collectedAt time.Time,
) (MarketItemSnapshot, error) {
	if marketName == "" {
		return MarketItemSnapshot{}, ErrEmptyMarketName
	}
	if externalItemID == "" {
		return MarketItemSnapshot{}, ErrEmptyExternalItemID
	}
	if marketHashName == "" {
		return MarketItemSnapshot{}, ErrEmptyMarketHashName
	}
	if len(price.Currency) != 3 {
		return MarketItemSnapshot{}, ErrInvalidCurrency
	}
	if price.AmountMinor <= 0 && event != EventItemRemoved {
		return MarketItemSnapshot{}, ErrInvalidSnapshotPrice
	}
	if availableCount < 0 {
		return MarketItemSnapshot{}, ErrNegativeListingCount
	}
	if event == "" {
		return MarketItemSnapshot{}, ErrEmptySnapshotEvent
	}
	if collectedAt.IsZero() {
		return MarketItemSnapshot{}, ErrEmptySnapshotTime
	}

	return MarketItemSnapshot{
		Market:         marketName,
		ExternalItemID: externalItemID,
		MarketHashName: marketHashName,
		Price:          price,
		AvailableCount: availableCount,
		ItemType:       itemType,
		Event:          event,
		CollectedAt:    collectedAt.UTC(),
	}, nil
}
