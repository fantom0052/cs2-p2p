package scheduler

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/fanto/p2p/steam-collector/internal/usecase/collector"
)

// Options содержит настройки цикла сбора данных.
type Options struct {
	Interval       time.Duration
	RateLimitPause time.Duration
	Jitter         time.Duration
	PageSize       int
	Query          string
}

// Scheduler периодически запускает сбор страниц Steam Market.
type Scheduler struct {
	collector Collector
	logger    *slog.Logger
	options   Options
}

// NewScheduler создает scheduler adapter.
func NewScheduler(collector Collector, logger *slog.Logger, options Options) (*Scheduler, error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &Scheduler{
		collector: collector,
		logger:    logger,
		options:   options,
	}, nil
}

// Run запускает цикл сбора данных до отмены context-а.
func (s *Scheduler) Run(ctx context.Context) error {
	start := 0

	for {
		page, err := s.collect(ctx, start)
		wait := s.nextWait(err)
		if err != nil {
			s.logger.ErrorContext(ctx, "collect page", "error", err, "start", start, "next_wait", wait)
		} else {
			s.logger.InfoContext(
				ctx,
				"page collected",
				"items", len(page.Items),
				"start", page.Start,
				"page_size", page.PageSize,
				"total_count", page.TotalCount,
			)
			start = nextStart(page)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// collect запускает одну итерацию сбора для конкретного offset-а Steam Market.
func (s *Scheduler) collect(ctx context.Context, start int) (collector.FetchPage, error) {
	return s.collector.CollectPage(ctx, collector.FetchParams{
		Query: s.options.Query,
		Start: start,
		Count: s.options.PageSize,
	})
}

func (s *Scheduler) nextWait(err error) time.Duration {
	wait := s.options.Interval
	if err != nil && errors.Is(err, collector.ErrRateLimited) {
		wait = s.options.RateLimitPause
	}

	return wait + randomJitter(s.options.Jitter)
}

func randomJitter(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}

	return time.Duration(rand.Int63n(int64(max)))
}

// nextStart вычисляет offset следующей страницы или возвращает 0, если текущий проход завершен.
func nextStart(page collector.FetchPage) int {
	if page.TotalCount <= 0 || page.PageSize <= 0 {
		return 0
	}

	next := page.Start + page.PageSize
	if next >= page.TotalCount {
		return 0
	}

	return next
}

// validateOptions проверяет настройки scheduler-а перед запуском.
func validateOptions(options Options) error {
	if options.Interval <= 0 {
		return ErrInvalidInterval
	}
	if options.RateLimitPause <= 0 {
		return ErrInvalidRateLimitPause
	}
	if options.Jitter < 0 {
		return ErrInvalidJitter
	}
	if options.PageSize <= 0 {
		return ErrInvalidPageSize
	}

	return nil
}
