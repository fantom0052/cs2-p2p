package steam

type searchRenderResponse struct {
	Success    bool               `json:"success"`
	Start      int                `json:"start"`
	PageSize   int                `json:"pagesize"`
	TotalCount int                `json:"total_count"`
	Results    []searchRenderItem `json:"results"`
}

type searchRenderItem struct {
	Name             string           `json:"name"`
	HashName         string           `json:"hash_name"`
	SellListings     int              `json:"sell_listings"`
	SellPrice        int              `json:"sell_price"`
	SellPriceText    string           `json:"sell_price_text"`
	AssetDescription assetDescription `json:"asset_description"`
	SalePriceText    string           `json:"sale_price_text"`
}

type assetDescription struct {
	AppID          int    `json:"appid"`
	ClassID        string `json:"classid"`
	IconURL        string `json:"icon_url"`
	Tradable       int    `json:"tradable"`
	Type           string `json:"type"`
	MarketName     string `json:"market_name"`
	MarketHashName string `json:"market_hash_name"`
	Commodity      int    `json:"commodity"`
}
