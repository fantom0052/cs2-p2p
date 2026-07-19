package collector

import (
	"context"

	"github.com/fanto/p2p/steam-collector/internal/domain/market"
)

// FetchParams описывает параметры одной итерации сбора данных.
type FetchParams struct {
	Query      string
	SortColumn string
	SortDir    string
	Start      int
	Count      int
}

// FetchPage содержит одну страницу собранных snapshot-ов.
type FetchPage struct {
	Items      []market.MarketItemSnapshot
	Start      int
	PageSize   int
	TotalCount int
}

// UseCase реализует сценарий сбора market item snapshot-ов и передачи их дальше.
type UseCase struct {
	fetcher   MarketItemFetcher
	publisher SnapshotPublisher
}

// NewUsecase создает collector usecase.
func NewUsecase(fetcher MarketItemFetcher, publisher SnapshotPublisher) *UseCase {
	return &UseCase{
		fetcher:   fetcher,
		publisher: publisher,
	}
}

// CollectPage собирает одну страницу snapshot-ов и публикует результат.
func (uc *UseCase) CollectPage(ctx context.Context, params FetchParams) (FetchPage, error) {
	page, err := uc.fetcher.FetchMarketItems(ctx, params)
	if err != nil {
		return FetchPage{}, err
	}

	if len(page.Items) == 0 {
		return page, nil
	}

	if err := uc.publisher.PublishSnapshots(ctx, page.Items); err != nil {
		return FetchPage{}, err
	}

	return page, nil
}
