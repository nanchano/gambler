package coingecko

import (
	"reflect"
	"testing"

	"github.com/nanchano/gambler/internal/core"
)

func TestBuildExtraField(t *testing.T) {
	r := newTestResponse()
	name := "Parse Image, MarketData, CommunityData, DeveloperData and PublicInterestStats into an aggregated map"
	expected := map[string]interface{}{
		"thumb":         "link_to_thumb",
		"current_price": map[string]float64{"usd": 123.45},
		"market_cap":    map[string]float64{"usd": 678.90},
		"total_volume":  map[string]float64{"usd": 432.10},
		"facebook":      map[string]int32{"likes": 1234},
		"stars":         123456,
		"alexa_rank":    3000,
	}
	t.Run(name, func(t *testing.T) {
		actual := r.buildExtraField()
		if !reflect.DeepEqual(actual, expected) {
			t.Fail()
		}
	})
}

func TestConvert(t *testing.T) {
	r := newTestResponse()
	name := "Convert the coingeckoResponse into a GamblerEvent"
	expected := core.GamblerEvent{
		ID:        r.ID,
		Name:      r.Name,
		Symbol:    r.Symbol,
		Date:      r.Date,
		Price:     r.MarketData["current_price"]["usd"],
		MarketCap: r.MarketData["market_cap"]["usd"],
		Volume:    r.MarketData["total_volume"]["usd"],
		Extra:     r.buildExtraField(),
	}
	t.Run(name, func(t *testing.T) {
		actual := r.Convert()
		if !reflect.DeepEqual(actual, &expected) {
			t.Fail()
		}
	})
}

func newTestResponse() *coingeckoResponse {
	return &coingeckoResponse{
		ID:     "ethereum",
		Symbol: "eth",
		Name:   "Ethereum",
		Image:  map[string]string{"thumb": "link_to_thumb"},
		MarketData: map[string]map[string]float64{
			"current_price": {"usd": 123.45},
			"market_cap":    {"usd": 678.90},
			"total_volume":  {"usd": 432.10},
		},
		CommunityData:       map[string]interface{}{"facebook": map[string]int32{"likes": 1234}},
		DeveloperData:       map[string]interface{}{"stars": 123456},
		PublicInterestStats: map[string]interface{}{"alexa_rank": 3000},
		Date:                "20-04-2022",
	}
}
