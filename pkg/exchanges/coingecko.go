package exchanges

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/nanchano/gambler/internal/core"
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

// CoingeckoPipeline handles API requests for a given coin ID
type CoingeckoPipeline struct {
	URL string
	ID  string
}

// NewCoingeckoPipeline creates a new CoingeckoPipeline with a default URL
func NewCoingeckoPipeline(id string) *CoingeckoPipeline {
	return &CoingeckoPipeline{
		URL: "https://api.coingecko.com/api/v3/",
		ID:  id,
	}
}

// Extract retrieves a response from the Coingecko API
func (c *CoingeckoPipeline) Extract(dates ...string) <-chan *coingeckoResponse {
	out := make(chan *coingeckoResponse)
	go func() {
		defer close(out)
		for _, date := range dates {
			url := c.prepareURL(date)
			log.Println(date)
			// log.Printf("Requesting: %s", url)
			resp, err := http.Get(url)

			if err != nil {
				log.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var cr coingeckoResponse
			if err := json.Unmarshal(body, &cr); err != nil {
				log.Fatal("Can not unmarshal JSON into struct")
			}

			cr.Date = date
			out <- &cr
		}
	}()
	return out
}

// Process normalizes the coingeckoResponse into a core.GamblerEvent
func (c *CoingeckoPipeline) Process(in <-chan *coingeckoResponse) <-chan *core.GamblerEvent {
	out := make(chan *core.GamblerEvent)
	go func() {
		defer close(out)
		for resp := range in {
			log.Println("Processing response")
			gr := processCoingeckoResponse(resp)
			out <- gr
		}
	}()
	return out
}

// Save dumps the event into a JSON file
func (c *CoingeckoPipeline) Save(in <-chan *core.GamblerEvent) {
	for {
		i, ok := <-in
		if ok {
			i.ToJSON()
		} else {
			return
		}
	}
}

// prepareURL prepares the URL (path + query params) for a request
func (c *CoingeckoPipeline) prepareURL(date string) string {
	base, err := url.Parse(c.URL)
	if err != nil {
		log.Fatal(err)
	}
	base.Path += fmt.Sprintf("coins/%s/history", c.ID)
	params := url.Values{}
	params.Add("date", date)
	params.Add("localization", "false")
	base.RawQuery = params.Encode()
	return base.String()
}

// processCoingeckoResponse transforms a coingeckoResponse into a core.GamblerEvent
func processCoingeckoResponse(cr *coingeckoResponse) *core.GamblerEvent {
	var gr core.GamblerEvent
	gr.ID = cr.ID
	gr.Name = cr.Name
	gr.Symbol = cr.Symbol
	gr.Date = cr.Date
	gr.Price = cr.MarketData["current_price"]["usd"]
	gr.MarketCap = cr.MarketData["market_cap"]["usd"]
	gr.Volume = cr.MarketData["total_volume"]["usd"]
	gr.Extra = cr.buildExtraField()

	return &gr
}

// buildExtraField appends the maps into a single one
func (cr *coingeckoResponse) buildExtraField() map[string]any {
	extra := make(map[string]any)
	for k, v := range cr.Image {
		extra[k] = v
	}
	for k, v := range cr.MarketData {
		extra[k] = v
	}
	for k, v := range cr.CommunityData {
		extra[k] = v
	}
	for k, v := range cr.DeveloperData {
		extra[k] = v
	}
	for k, v := range cr.PublicInterestStats {
		extra[k] = v
	}
	return extra
}
