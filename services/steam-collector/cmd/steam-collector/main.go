package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kafkaAdapter "github.com/fanto/p2p/steam-collector/internal/adapters/kafka"
	schedulerAdapter "github.com/fanto/p2p/steam-collector/internal/adapters/scheduler"
	steamAdapter "github.com/fanto/p2p/steam-collector/internal/adapters/steam"
	"github.com/fanto/p2p/steam-collector/internal/config"
	"github.com/fanto/p2p/steam-collector/internal/domain/market"
	"github.com/fanto/p2p/steam-collector/internal/usecase/collector"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	steamOptions, err := newSteamOptions(cfg.Steam)
	if err != nil {
		logger.Error("build steam options", "error", err)
		os.Exit(1)
	}

	steamClient, err := steamAdapter.NewClient(&http.Client{Timeout: cfg.Steam.RequestTimeout}, steamOptions)
	if err != nil {
		logger.Error("create steam client", "error", err)
		os.Exit(1)
	}

	publisher, err := kafkaAdapter.NewPublisher(ctx, newKafkaOptions(cfg.Kafka))
	if err != nil {
		logger.Error("create kafka publisher", "error", err)
		os.Exit(1)
	}
	defer publisher.Close()

	collectorUsecase := collector.NewUsecase(steamClient, publisher)

	scheduler, err := schedulerAdapter.NewScheduler(collectorUsecase, logger, newSchedulerOptions(cfg.Scheduler))
	if err != nil {
		logger.Error("create scheduler", "error", err)
		os.Exit(1)
	}

	logger.Info("steam collector started")
	if err := scheduler.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("run scheduler", "error", err)
		os.Exit(1)
	}
	logger.Info("steam collector stopped")
}

func newSteamOptions(cfg config.Steam) (steamAdapter.ClientOptions, error) {
	currency, err := market.NewCurrency(cfg.Currency)
	if err != nil {
		return steamAdapter.ClientOptions{}, err
	}

	return steamAdapter.ClientOptions{
		HomeURL:            cfg.HomeURL,
		BaseURL:            cfg.BaseURL,
		ImageBaseURL:       cfg.ImageBaseURL,
		AppID:              cfg.AppID,
		CurrencyID:         cfg.CurrencyID,
		Currency:           currency,
		UserAgent:          cfg.UserAgent,
		AcceptLanguage:     cfg.AcceptLanguage,
		CookieSettings:     cfg.CookieSettings,
		BootstrapCookies:   cfg.BootstrapCookies,
		RequestTimeout:     cfg.RequestTimeout,
		NoRender:           cfg.NoRender,
		DefaultSortColumn:  cfg.DefaultSortColumn,
		DefaultSortDir:     cfg.DefaultSortDir,
		DefaultPageSize:    cfg.DefaultPageSize,
		SearchDescriptions: cfg.SearchDescriptions,
		CategoryFilters: map[string]string{
			"category_730_ItemSet[]":        cfg.CategoryItemSet,
			"category_730_ProPlayer[]":      cfg.CategoryProPlayer,
			"category_730_StickerCapsule[]": cfg.CategoryStickerCapsule,
			"category_730_Tournament[]":     cfg.CategoryTournament,
			"category_730_TournamentTeam[]": cfg.CategoryTournamentTeam,
			"category_730_Type[]":           cfg.CategoryType,
			"category_730_Weapon[]":         cfg.CategoryWeapon,
		},
	}, nil
}

func newKafkaOptions(cfg config.Kafka) kafkaAdapter.PublisherOptions {
	return kafkaAdapter.PublisherOptions{
		Brokers:       cfg.Brokers,
		Topic:         cfg.Topic,
		ClientID:      cfg.ClientID,
		Acks:          cfg.Acks,
		Compression:   cfg.Compression,
		Linger:        cfg.Linger,
		BatchBytes:    cfg.BatchBytes,
		RecordRetries: cfg.RecordRetries,
	}
}

func newSchedulerOptions(cfg config.Scheduler) schedulerAdapter.Options {
	return schedulerAdapter.Options{
		Interval:       cfg.Interval,
		RateLimitPause: cfg.RateLimitPause,
		Jitter:         cfg.Jitter,
		PageSize:       cfg.PageSize,
		Query:          cfg.Query,
	}
}
