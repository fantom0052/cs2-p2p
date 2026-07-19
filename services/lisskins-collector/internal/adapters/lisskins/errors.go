package lisskins

import "errors"

var (
	ErrEmptyDumpURL      = errors.New("empty lisskins dump url")
	ErrEmptyWebSocketURL = errors.New("empty lisskins websocket url")
	ErrEmptyTokenURL     = errors.New("empty lisskins token url")
	ErrEmptyAPIKey       = errors.New("empty lisskins api key")
	ErrInvalidBatchSize  = errors.New("invalid lisskins batch size")
)
