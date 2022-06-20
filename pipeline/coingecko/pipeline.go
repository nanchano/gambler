package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/nanchano/gambler/internal/core"
)

// Pipeline handles API requests for a given coin ID
type Pipeline struct {
	URL string
	ID  string
}

// NewPipeline creates a new Pipeline with a default URL
func NewPipeline(id string) *Pipeline {
	return &Pipeline{
		URL: "https://api.coingecko.com/api/v3/",
		ID:  id,
	}
}

// Extract retrieves the response from the Coingecko API for the given dates
func (p *Pipeline) Extract(dates ...string) <-chan core.PipelineResponse {
	out := make(chan core.PipelineResponse)
	go func() {
		defer close(out)
		for _, date := range dates {
			url := p.prepareURL(date)
			log.Printf("Requesting: %s", url)
			resp, err := http.Get(url)

			if err != nil {
				log.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var response coingeckoResponse
			if err := json.Unmarshal(body, &response); err != nil {
				log.Fatalf("Could not unmarshal JSON into struct: %s", err)
			}

			response.Date = date
			out <- &response
		}
	}()
	return out
}

// Process normalizes the coingecko responses into a core.GamblerEvents
func (p *Pipeline) Process(responses <-chan core.PipelineResponse) <-chan *core.GamblerEvent {
	out := make(chan *core.GamblerEvent)
	go func() {
		defer close(out)
		for resp := range responses {
			log.Println("Transforming the Coingecko Response into a GamblerEvent")
			event := resp.Convert()
			out <- event
		}
	}()
	return out
}

// prepareURL prepares the URL (path + query params) for a request
func (p *Pipeline) prepareURL(date string) string {
	base, err := url.Parse(p.URL)
	if err != nil {
		log.Fatal(err)
	}
	base.Path += fmt.Sprintf("coins/%s/history", p.ID)
	params := url.Values{}
	params.Add("date", date)
	params.Add("localization", "false")
	base.RawQuery = params.Encode()
	return base.String()
}
