package lisskins

type dumpResponse struct {
	Items      []dumpItem `json:"items"`
	LastUpdate uint64     `json:"last_update"`
	Status     string     `json:"status"`
}

type dumpItem struct {
	ID             int     `json:"id"`
	CreatedAt      string  `json:"created_at"`
	GameID         int     `json:"game_id"`
	ItemAssetID    string  `json:"item_asset_id"`
	ItemClassID    string  `json:"item_class_id"`
	ItemFloat      *string `json:"item_float"`
	ItemPaintIndex *int64  `json:"item_paint_index"`
	ItemPaintSeed  *int64  `json:"item_paint_seed"`
	Name           string  `json:"name"`
	NameTag        *string `json:"name_tag"`
	Price          float64 `json:"price"`
	UnlockAt       *string `json:"unlock_at"`
	FloatName      *string `json:"float_name"`
}

type skinEvent struct {
	CreatedAt      string  `json:"created_at"`
	Event          string  `json:"event"`
	GameID         int     `json:"game_id"`
	ID             int     `json:"id"`
	ItemAssetID    string  `json:"item_asset_id"`
	ItemClassID    string  `json:"item_class_id"`
	ItemPaintIndex *int    `json:"item_paint_index,omitempty"`
	ItemPaintSeed  *int    `json:"item_paint_seed,omitempty"`
	Name           string  `json:"name"`
	NameTag        *string `json:"name_tag,omitempty"`
	Price          float64 `json:"price"`
	UnlockAt       string  `json:"unlock_at"`
}

type tokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}
