package steam

import "errors"

var (
	ErrSearchResponseFailed      = errors.New("steam search response is not successful")
	ErrEmptyHomeURL              = errors.New("steam home url is empty")
	ErrInvalidHomeURL            = errors.New("steam home url must end with slash")
	ErrEmptyBaseURL              = errors.New("steam base url is empty")
	ErrEmptyImageBaseURL         = errors.New("steam image base url is empty")
	ErrInvalidAppID              = errors.New("steam app id must be positive")
	ErrInvalidCurrencyID         = errors.New("steam currency id must be positive")
	ErrEmptyUserAgent            = errors.New("steam user agent is empty")
	ErrEmptyAcceptLanguage       = errors.New("steam accept language is empty")
	ErrEmptyDefaultSortColumn    = errors.New("steam default sort column is empty")
	ErrEmptyDefaultSortDir       = errors.New("steam default sort dir is empty")
	ErrInvalidDefaultPageSize    = errors.New("steam default page size must be positive")
	ErrInvalidSearchDescriptions = errors.New("steam search descriptions must be 0 or 1")
	ErrInvalidNoRender           = errors.New("steam norender must be 0 or 1")
	ErrInvalidRequestTimeout     = errors.New("steam request timeout must be positive")
)
