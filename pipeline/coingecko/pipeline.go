package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/nanchano/gambler/core"
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

// Run extracts coingeckoResponses, processes them into core.GamblerEvents
// and stores them on a given core.Repository for a set of dates
func (p *Pipeline) Run(repo core.GamblerRepository, dates ...string) {
	responses := p.extract(dates...)
	events := p.process(responses)
	p.store(events, repo)
}

// extract retrieves the response from the Coingecko API for the given dates
func (p *Pipeline) extract(dates ...string) <-chan core.PipelineResponse {
	out := make(chan core.PipelineResponse)
	go func() {
		defer close(out)
		fmt.Println("ASDASDAS")
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

// process normalizes the coingeckoResponses into core.GamblerEvents
func (p *Pipeline) process(responses <-chan core.PipelineResponse) <-chan *core.GamblerEvent {
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

// store saves the relevant core.GamblerEvents on a given core.Repository
func (p *Pipeline) store(events <-chan *core.GamblerEvent, repo core.GamblerRepository) {
	for event := range events {
		repo.Store(event)
	}
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
