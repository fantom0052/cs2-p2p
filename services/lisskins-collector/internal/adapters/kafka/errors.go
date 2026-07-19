package kafka

import "errors"

var (
	ErrEmptyBrokers        = errors.New("empty kafka brokers")
	ErrEmptyTopic          = errors.New("empty kafka topic")
	ErrEmptyClientID       = errors.New("empty kafka client id")
	ErrInvalidAcks         = errors.New("invalid kafka acks")
	ErrInvalidCompression  = errors.New("invalid kafka compression")
	ErrInvalidBatchTimeout = errors.New("invalid kafka batch timeout")
	ErrInvalidBatchBytes   = errors.New("invalid kafka batch bytes")
)
