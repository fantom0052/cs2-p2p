package collector

import (
	"context"

	"github.com/fanto/p2p/steam-collector/internal/domain/market"
)

// MarketItemFetcher описывает источник рыночных данных для collector usecase.
type MarketItemFetcher interface {
	FetchMarketItems(ctx context.Context, params FetchParams) (FetchPage, error)
}

// SnapshotPublisher описывает выходной порт для публикации собранных snapshot-ов.
type SnapshotPublisher interface {
	PublishSnapshots(ctx context.Context, snapshots []market.MarketItemSnapshot) error
}
