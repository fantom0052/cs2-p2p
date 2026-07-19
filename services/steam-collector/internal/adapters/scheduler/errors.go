package scheduler

import "errors"

var (
	ErrInvalidInterval       = errors.New("scheduler interval must be positive")
	ErrInvalidRateLimitPause = errors.New("scheduler rate limit pause must be positive")
	ErrInvalidJitter         = errors.New("scheduler jitter must be non-negative")
	ErrInvalidPageSize       = errors.New("scheduler page size must be positive")
)
