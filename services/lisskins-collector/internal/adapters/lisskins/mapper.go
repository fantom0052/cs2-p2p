package lisskins

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
)

func mapDumpItem(item dumpItem, currency market.Currency, collectedAt time.Time) (market.MarketItemSnapshot, error) {
	return newSnapshot(strconv.Itoa(item.ID), item.Name, item.Price, currency, market.EventSnapshot, collectedAt)
}

func mapSkinEvent(event skinEvent, currency market.Currency, collectedAt time.Time) (market.MarketItemSnapshot, error) {
	return newSnapshot(strconv.Itoa(event.ID), event.Name, event.Price, currency, mapEventName(event.Event), collectedAt)
}

func newSnapshot(
	externalID string,
	name string,
	price float64,
	currency market.Currency,
	event market.SnapshotEvent,
	collectedAt time.Time,
) (market.MarketItemSnapshot, error) {
	externalItemID, err := market.NewExternalItemID(externalID)
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	marketHashName, err := market.NewMarketHashName(name)
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	money, err := market.NewMoney(int64(math.Round(price*100)), currency)
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	return market.NewMarketItemSnapshot(
		market.MarketLisSkins,
		externalItemID,
		marketHashName,
		money,
		1,
		detectItemType(name),
		event,
		collectedAt,
	)
}

func mapEventName(value string) market.SnapshotEvent {
	switch value {
	case "obtained_skin_added":
		return market.EventItemAdded
	case "obtained_skin_price_changed":
		return market.EventPriceChanged
	case "obtained_skin_deleted":
		return market.EventItemRemoved
	default:
		return market.EventSnapshot
	}
}

func detectItemType(name string) market.ItemType {
	name = strings.TrimSpace(name)
	if name == "" {
		return "unknown"
	}

	if strings.HasPrefix(name, "Sticker |") {
		return "sticker"
	}
	if strings.Contains(name, "Case") {
		return "case"
	}
	if strings.HasPrefix(name, "★") {
		return "knife_or_gloves"
	}
	if idx := strings.Index(name, "|"); idx > 0 {
		return market.ItemType(strings.ToLower(strings.TrimSpace(name[:idx])))
	}

	return "unknown"
}
