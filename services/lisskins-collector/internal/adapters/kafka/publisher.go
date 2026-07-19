package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/fanto/p2p/lisskins-collector/internal/domain/market"
)

// Publisher публикует собранные snapshot-ы в Kafka-compatible broker.
type Publisher struct {
	options  PublisherOptions
	producer *kgo.Client
}

// PublisherOptions содержит настройки Kafka publisher-а.
type PublisherOptions struct {
	Brokers       []string
	Topic         string
	ClientID      string
	Acks          string
	Compression   string
	Linger        time.Duration
	BatchBytes    int32
	RecordRetries int
}

// NewPublisher создает Kafka adapter для публикации snapshot-ов.
func NewPublisher(ctx context.Context, options PublisherOptions) (*Publisher, error) {
	if err := validatePublisherOptions(options); err != nil {
		return nil, err
	}

	clientOptions, err := newClientOptions(options)
	if err != nil {
		return nil, err
	}

	producer, err := kgo.NewClient(clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	if err := producer.EnsureProduceConnectionIsOpen(ctx, -1); err != nil {
		producer.Close()
		return nil, fmt.Errorf("ensure kafka produce connection is open: %w", err)
	}

	return &Publisher{
		options:  options,
		producer: producer,
	}, nil
}

// PublishSnapshots публикует batch snapshot-ов в Kafka.
func (p *Publisher) PublishSnapshots(ctx context.Context, snapshots []market.MarketItemSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}

	if err := p.producer.EnsureProduceConnectionIsOpen(ctx, -1); err != nil {
		return fmt.Errorf("ensure kafka produce connection is open: %w", err)
	}

	records := make([]*kgo.Record, 0, len(snapshots))
	for _, snapshot := range snapshots {
		record, err := newRecord(p.options.Topic, snapshot)
		if err != nil {
			return err
		}
		records = append(records, record)
	}

	if err := p.producer.ProduceSync(ctx, records...).FirstErr(); err != nil {
		return fmt.Errorf("send kafka records: %w", err)
	}

	return nil
}

// Close закрывает Kafka producer и освобождает сетевые ресурсы.
func (p *Publisher) Close() {
	if p.producer == nil {
		return
	}

	p.producer.Close()
}

func newRecord(topic string, snapshot market.MarketItemSnapshot) (*kgo.Record, error) {
	payload, err := json.Marshal(newSnapshotMessage(snapshot))
	if err != nil {
		return nil, err
	}

	return &kgo.Record{
		Topic:     topic,
		Key:       []byte(snapshot.ExternalItemID),
		Value:     payload,
		Timestamp: snapshot.CollectedAt,
	}, nil
}

func newClientOptions(options PublisherOptions) ([]kgo.Opt, error) {
	acks, err := parseAcks(options.Acks)
	if err != nil {
		return nil, err
	}

	compression, err := parseCompression(options.Compression)
	if err != nil {
		return nil, err
	}

	return []kgo.Opt{
		kgo.SeedBrokers(options.Brokers...),
		kgo.ClientID(options.ClientID),
		kgo.RequiredAcks(acks),
		kgo.ProducerBatchCompression(compression),
		kgo.ProducerLinger(options.Linger),
		kgo.ProducerBatchMaxBytes(options.BatchBytes),
		kgo.RecordRetries(options.RecordRetries),
	}, nil
}

func validatePublisherOptions(options PublisherOptions) error {
	if len(options.Brokers) == 0 {
		return ErrEmptyBrokers
	}
	if options.Topic == "" {
		return ErrEmptyTopic
	}
	if options.ClientID == "" {
		return ErrEmptyClientID
	}
	if options.Linger <= 0 {
		return ErrInvalidBatchTimeout
	}
	if options.BatchBytes <= 0 {
		return ErrInvalidBatchBytes
	}

	return nil
}

func parseAcks(value string) (kgo.Acks, error) {
	switch value {
	case "all":
		return kgo.AllISRAcks(), nil
	case "leader":
		return kgo.LeaderAck(), nil
	case "none":
		return kgo.NoAck(), nil
	default:
		return kgo.Acks{}, ErrInvalidAcks
	}
}

func parseCompression(value string) (kgo.CompressionCodec, error) {
	switch value {
	case "none":
		return kgo.NoCompression(), nil
	case "gzip":
		return kgo.GzipCompression(), nil
	case "snappy":
		return kgo.SnappyCompression(), nil
	case "lz4":
		return kgo.Lz4Compression(), nil
	case "zstd":
		return kgo.ZstdCompression(), nil
	default:
		return kgo.CompressionCodec{}, ErrInvalidCompression
	}
}
