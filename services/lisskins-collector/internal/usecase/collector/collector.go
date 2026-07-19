package collector

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
)

// UseCase реализует фоновый сценарий сбора Lis-Skins: начальный dump и live-события.
type UseCase struct {
	source    SnapshotSource
	publisher SnapshotPublisher
	logger    *slog.Logger
	batchSize int
}

// NewUsecase создает collector usecase.
func NewUsecase(source SnapshotSource, publisher SnapshotPublisher, logger *slog.Logger, batchSize int) (*UseCase, error) {
	if batchSize <= 0 {
		return nil, ErrEmptyBatchSize
	}

	return &UseCase{
		source:    source,
		publisher: publisher,
		logger:    logger,
		batchSize: batchSize,
	}, nil
}

// Run загружает текущий dump Lis-Skins, публикует его batch-ами и затем слушает WebSocket.
func (uc *UseCase) Run(ctx context.Context) error {
	snapshots, err := uc.source.LoadDump(ctx)
	if err != nil {
		return fmt.Errorf("load dump: %w", err)
	}

	if err := uc.publishBatches(ctx, snapshots); err != nil {
		return fmt.Errorf("publish dump: %w", err)
	}

	uc.logger.Info("lisskins dump published", "items", len(snapshots))

	return uc.source.Watch(ctx, func(ctx context.Context, snapshot market.MarketItemSnapshot) error {
		return uc.publisher.PublishSnapshots(ctx, []market.MarketItemSnapshot{snapshot})
	})
}

func (uc *UseCase) publishBatches(ctx context.Context, snapshots []market.MarketItemSnapshot) error {
	for i := 0; i < len(snapshots); i += uc.batchSize {
		end := i + uc.batchSize
		if end > len(snapshots) {
			end = len(snapshots)
		}

		if err := uc.publisher.PublishSnapshots(ctx, snapshots[i:end]); err != nil {
			return err
		}

		uc.logger.Info("lisskins dump batch published", "from", i, "to", end)
	}

	return nil
}
