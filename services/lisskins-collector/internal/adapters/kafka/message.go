package kafka

import (
	"time"

	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
)

type snapshotMessage struct {
	Market           string    `json:"market"`
	ExternalItemID   string    `json:"external_item_id"`
	MarketHashName   string    `json:"market_hash_name"`
	PriceAmountMinor int64     `json:"price_amount_minor"`
	Currency         string    `json:"currency"`
	AvailableCount   int       `json:"available_count"`
	ItemType         string    `json:"item_type"`
	Event            string    `json:"event"`
	CollectedAt      time.Time `json:"collected_at"`
}

func newSnapshotMessage(snapshot market.MarketItemSnapshot) snapshotMessage {
	return snapshotMessage{
		Market:           string(snapshot.Market),
		ExternalItemID:   string(snapshot.ExternalItemID),
		MarketHashName:   string(snapshot.MarketHashName),
		PriceAmountMinor: snapshot.Price.AmountMinor,
		Currency:         string(snapshot.Price.Currency),
		AvailableCount:   snapshot.AvailableCount,
		ItemType:         string(snapshot.ItemType),
		Event:            string(snapshot.Event),
		CollectedAt:      snapshot.CollectedAt,
	}
}
