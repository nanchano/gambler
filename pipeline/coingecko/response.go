package coingecko

import (
	"github.com/nanchano/gambler/core"
)

type coingeckoResponse struct {
	ID                  string                        `json:"id"`
	Symbol              string                        `json:"symbol"`
	Name                string                        `json:"name"`
	Image               map[string]string             `json:"image"`
	MarketData          map[string]map[string]float64 `json:"market_data"`
	CommunityData       map[string]interface{}        `json:"community_data"`
	DeveloperData       map[string]interface{}        `json:"developer_data"`
	PublicInterestStats map[string]interface{}        `json:"public_interest_stats"`
	Date                string
}

// processCoingeckoResponse transforms a coingeckoResponse into a core.GamblerEvent
func (response *coingeckoResponse) Convert() *core.GamblerEvent {
	var event core.GamblerEvent
	event.ID = response.ID
	event.Name = response.Name
	event.Symbol = response.Symbol
	event.Date = response.Date
	event.Price = response.MarketData["current_price"]["usd"]
	event.MarketCap = response.MarketData["market_cap"]["usd"]
	event.Volume = response.MarketData["total_volume"]["usd"]
	event.Extra = response.buildExtraField()

	return &event
}

// buildExtraField appends the maps into a single one
func (response *coingeckoResponse) buildExtraField() map[string]interface{} {
	extra := make(map[string]interface{})
	for k, v := range response.Image {
		extra[k] = v
	}
	for k, v := range response.MarketData {
		extra[k] = v
	}
	for k, v := range response.CommunityData {
		extra[k] = v
	}
	for k, v := range response.DeveloperData {
		extra[k] = v
	}
	for k, v := range response.PublicInterestStats {
		extra[k] = v
	}
	return extra
}
