package core

// GamblerEvent represents a normalized response from any crypto API
type GamblerEvent struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Symbol    string      `json:"symbol"`
	Date      string      `json:"date"`
	Price     float64     `json:"price"`
	MarketCap float64     `json:"market_cap"`
	Volume    float64     `json:"volume"`
	Extra     interface{} `json:"extra"`
}
