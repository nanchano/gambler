package coingecko

import (
	"reflect"
	"testing"

	"github.com/nanchano/gambler/internal/core"
)

func TestBuildExtraField(t *testing.T) {
	cases := []struct {
		name     string
		response *coingeckoResponse
		expected map[string]interface{}
	}{
		{
			name:     "Has data",
			response: newTestResponse(),
			expected: map[string]interface{}{
				"thumb":         "link_to_thumb",
				"current_price": map[string]float64{"usd": 123.45},
				"market_cap":    map[string]float64{"usd": 678.90},
				"total_volume":  map[string]float64{"usd": 432.10},
				"facebook":      map[string]int32{"likes": 1234},
				"stars":         123456,
				"alexa_rank":    3000,
			},
		},
		{
			name:     "Empty response",
			response: &coingeckoResponse{},
			expected: map[string]interface{}{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.response.buildExtraField()
			if !reflect.DeepEqual(actual, c.expected) {
				t.Fail()
			}
		})
	}
}

func TestConvert(t *testing.T) {
	r := newTestResponse()
	cases := []struct {
		name     string
		response *coingeckoResponse
		expected *core.GamblerEvent
	}{
		{
			name:     "Has data",
			response: r,
			expected: &core.GamblerEvent{
				ID:        r.ID,
				Name:      r.Name,
				Symbol:    r.Symbol,
				Date:      r.Date,
				Price:     r.MarketData["current_price"]["usd"],
				MarketCap: r.MarketData["market_cap"]["usd"],
				Volume:    r.MarketData["total_volume"]["usd"],
				Extra:     r.buildExtraField(),
			},
		},
		{
			name:     "Empty response",
			response: &coingeckoResponse{},
			expected: &core.GamblerEvent{
				Extra: map[string]interface{}{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.response.Convert()
			if !reflect.DeepEqual(actual, c.expected) {
				t.Fail()
			}
		})
	}
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
