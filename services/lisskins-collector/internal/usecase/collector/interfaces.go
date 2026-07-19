package collector

import (
	"context"

	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
)

// SnapshotHandler обрабатывает один snapshot, полученный из live-источника.
type SnapshotHandler func(ctx context.Context, snapshot market.MarketItemSnapshot) error

// SnapshotSource описывает источник dump-а и live-событий Lis-Skins.
type SnapshotSource interface {
	LoadDump(ctx context.Context) ([]market.MarketItemSnapshot, error)
	Watch(ctx context.Context, handler SnapshotHandler) error
}

// SnapshotPublisher описывает выходной порт для публикации snapshot-ов.
type SnapshotPublisher interface {
	PublishSnapshots(ctx context.Context, snapshots []market.MarketItemSnapshot) error
}
