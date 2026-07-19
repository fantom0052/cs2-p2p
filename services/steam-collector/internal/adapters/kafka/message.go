package kafka

import (
	"time"

	"github.com/fanto/p2p/steam-collector/internal/domain/market"
)

type snapshotMessage struct {
	Market           string    `json:"market"`
	MarketHashName   string    `json:"market_hash_name"`
	PriceAmountMinor int64     `json:"price_amount_minor"`
	Currency         string    `json:"currency"`
	AvailableCount   int       `json:"available_count"`
	ItemType         string    `json:"item_type"`
	ImageURL         string    `json:"image_url"`
	CollectedAt      time.Time `json:"collected_at"`
}

func newSnapshotMessage(snapshot market.MarketItemSnapshot) snapshotMessage {
	return snapshotMessage{
		Market:           string(snapshot.Market),
		MarketHashName:   string(snapshot.MarketHashName),
		PriceAmountMinor: snapshot.Price.AmountMinor,
		Currency:         string(snapshot.Price.Currency),
		AvailableCount:   snapshot.AvailableCount,
		ItemType:         string(snapshot.ItemType),
		ImageURL:         string(snapshot.ImageURL),
		CollectedAt:      snapshot.CollectedAt,
	}
}
