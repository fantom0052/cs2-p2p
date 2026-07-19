package kafka

import "errors"

var (
	ErrEmptyBrokers        = errors.New("kafka brokers are empty")
	ErrEmptyTopic          = errors.New("kafka topic is empty")
	ErrEmptyClientID       = errors.New("kafka client id is empty")
	ErrInvalidBatchTimeout = errors.New("kafka batch timeout must be positive")
	ErrInvalidBatchBytes   = errors.New("kafka batch bytes must be positive")
	ErrInvalidAcks         = errors.New("kafka acks value is invalid")
	ErrInvalidCompression  = errors.New("kafka compression value is invalid")
)
