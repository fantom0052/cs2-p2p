package market

import "errors"

var (
	ErrEmptyMarketName      = errors.New("market name is empty")
	ErrEmptyMarketHashName  = errors.New("market hash name is empty")
	ErrInvalidCurrency      = errors.New("currency must be ISO 4217 code")
	ErrInvalidMoneyAmount   = errors.New("money amount cannot be negative")
	ErrInvalidSnapshotPrice = errors.New("snapshot price must be positive")
	ErrEmptyItemType        = errors.New("item type is empty")
	ErrNegativeListingCount = errors.New("available count cannot be negative")
	ErrEmptyImageURL        = errors.New("image url is empty")
	ErrEmptySnapshotTime    = errors.New("collected at is empty")
)
