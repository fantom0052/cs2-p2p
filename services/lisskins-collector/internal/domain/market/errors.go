package market

import "errors"

var (
	ErrEmptyMarketName      = errors.New("empty market name")
	ErrEmptyExternalItemID  = errors.New("empty external item id")
	ErrEmptyMarketHashName  = errors.New("empty market hash name")
	ErrInvalidCurrency      = errors.New("invalid currency")
	ErrInvalidMoneyAmount   = errors.New("invalid money amount")
	ErrInvalidSnapshotPrice = errors.New("invalid snapshot price")
	ErrNegativeListingCount = errors.New("negative listing count")
	ErrEmptySnapshotTime    = errors.New("empty snapshot time")
	ErrEmptySnapshotEvent   = errors.New("empty snapshot event")
)
