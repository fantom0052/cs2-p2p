package steam

import (
	"strings"
	"time"

	"github.com/fanto/p2p/steam-collector/internal/domain/market"
)

func mapSearchRenderItem(item searchRenderItem, currency market.Currency, imageBaseURL string, collectedAt time.Time) (market.MarketItemSnapshot, error) {
	marketHashNameRaw := item.AssetDescription.MarketHashName
	if strings.TrimSpace(marketHashNameRaw) == "" {
		marketHashNameRaw = item.HashName
	}

	marketHashName, err := market.NewMarketHashName(marketHashNameRaw)
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	price, err := market.NewMoney(int64(item.SellPrice), currency)
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	itemType, err := market.NewItemType(item.AssetDescription.Type)
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	imageURL, err := market.NewImageURL(buildSteamImageURL(imageBaseURL, item.AssetDescription.IconURL))
	if err != nil {
		return market.MarketItemSnapshot{}, err
	}

	return market.NewMarketItemSnapshot(
		market.MarketSteam,
		marketHashName,
		price,
		item.SellListings,
		itemType,
		imageURL,
		collectedAt,
	)
}

func buildSteamImageURL(imageBaseURL string, iconURL string) string {
	iconURL = strings.TrimSpace(iconURL)
	if iconURL == "" {
		return ""
	}
	if strings.HasPrefix(iconURL, "http://") || strings.HasPrefix(iconURL, "https://") {
		return iconURL
	}

	return strings.TrimRight(imageBaseURL, "/") + "/" + iconURL
}
