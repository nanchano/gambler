package exchanges

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type CoingeckoResponse map[string]interface{}

// CoingeckoHandler handles API requests for a given coin ID and date
type CoingeckoHandler struct {
	URL  string
	ID   string
	Date string
}

// NewCoingeckoHandler creates a new CoingeckoHandler with a default URL
func NewCoingeckoHandler(id, date string) *CoingeckoHandler {
	return &CoingeckoHandler{
		URL:  "https://api.coingecko.com/api/v3/",
		ID:   id,
		Date: date,
	}
}

func (c *CoingeckoHandler) prepareUrl() string {
	base, err := url.Parse(c.URL)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	base.Path += fmt.Sprintf("coins/%s/history", c.ID)

	params := url.Values{}
	params.Add("date", c.Date)
	base.RawQuery = params.Encode()

	return base.String()
}

func (c *CoingeckoHandler) Extract() []byte {
	url := c.prepareUrl()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func (c *CoingeckoHandler) Transform(resp []byte) *CoingeckoResponse {
	var cr CoingeckoResponse
	// var cr map[string]interface{}
	json.Unmarshal(resp, &cr)
	return &cr
}
