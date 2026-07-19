package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kafkaAdapter "github.com/fanto/p2p/lisskins-collector/internal/adapters/kafka"
	lisskinsAdapter "github.com/fanto/p2p/lisskins-collector/internal/adapters/lisskins"
	"github.com/fanto/p2p/lisskins-collector/internal/config"
	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
	"github.com/fanto/p2p/lisskins-collector/internal/usecase/collector"
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

	lisskinsOptions, err := newLisSkinsOptions(cfg.LisSkins)
	if err != nil {
		logger.Error("build lisskins options", "error", err)
		os.Exit(1)
	}

	lisskinsClient, err := lisskinsAdapter.NewClient(&http.Client{Timeout: cfg.LisSkins.RequestTimeout}, lisskinsOptions)
	if err != nil {
		logger.Error("create lisskins client", "error", err)
		os.Exit(1)
	}

	publisher, err := kafkaAdapter.NewPublisher(ctx, newKafkaOptions(cfg.Kafka))
	if err != nil {
		logger.Error("create kafka publisher", "error", err)
		os.Exit(1)
	}
	defer publisher.Close()

	collectorUsecase, err := collector.NewUsecase(lisskinsClient, publisher, logger, cfg.LisSkins.BatchSize)
	if err != nil {
		logger.Error("create collector usecase", "error", err)
		os.Exit(1)
	}

	logger.Info("lisskins collector started")
	if err := collectorUsecase.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("run collector", "error", err)
		os.Exit(1)
	}
	logger.Info("lisskins collector stopped")
}

func newLisSkinsOptions(cfg config.LisSkins) (lisskinsAdapter.ClientOptions, error) {
	currency, err := market.NewCurrency(cfg.Currency)
	if err != nil {
		return lisskinsAdapter.ClientOptions{}, err
	}

	return lisskinsAdapter.ClientOptions{
		APIKey:       cfg.APIKey,
		DumpURL:      cfg.DumpURL,
		WebSocketURL: cfg.WebSocketURL,
		WSTokenURL:   cfg.WSTokenURL,
		Currency:     currency,
		DumpLimit:    cfg.DumpLimit,
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
