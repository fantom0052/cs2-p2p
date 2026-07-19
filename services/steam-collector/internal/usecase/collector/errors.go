package collector

import "errors"

// ErrRateLimited означает, что источник данных временно ограничил частоту запросов.
var ErrRateLimited = errors.New("market source rate limited")
