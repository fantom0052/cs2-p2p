package scheduler

import (
	"context"

	"github.com/fanto/p2p/steam-collector/internal/usecase/collector"
)

// Collector описывает usecase, который scheduler запускает по расписанию.
type Collector interface {
	CollectPage(ctx context.Context, params collector.FetchParams) (collector.FetchPage, error)
}
