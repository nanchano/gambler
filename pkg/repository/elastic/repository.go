package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/nanchano/gambler/internal/core"
)

// repository will implement the core.GamblerRepository interface for the Elastic Repository
type repository struct {
	client *elastic.Client
}

// newElasticClient returns an ElasticSearch client and validates connection through a Ping
func newElasticClient() (*elastic.Client, error) {
	client, err := elastic.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	// Quick ping validation
	res, err := client.Info()
	if err != nil {
		return nil, err
	}
	io.Copy(ioutil.Discard, res.Body)

	return client, nil
}

// NewRepository creates a new repository that implements core.GamblerRepository
func NewRepository() core.GamblerRepository {
	client, err := newElasticClient()
	if err != nil {
		log.Fatalf("Error initializing the ElasticSearch client: %s", err)
	}

	return &repository{
		client: client,
	}
}

// Find returns a core.GamblerEvent from the ElasticSearch repository based on the coin and date provided
func (r *repository) Find(coin, date string) (*core.GamblerEvent, error) {
	var event core.GamblerEvent
	ctx := context.Background()
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": fmt.Sprintf("(id:\"%s\") AND (date:\"%s\")", coin, date),
			},
		},
	}
	_ = json.NewEncoder(&buf).Encode(query)

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("coins"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error getting the response for %s - %s: %s", coin, date, err)
		return nil, err
	}
	defer res.Body.Close()

	var response elasticSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return nil, err
	}

	if len(response.Hits.Hits) == 1 {
		if err := json.Unmarshal(response.Hits.Hits[0].Source, &event); err != nil {
			return nil, err
		}
		return &event, nil
	}
	return nil, errors.New("More than one result for the given filters")

}

// Store saves the core.GamblerEvents provided into the ElasticSearch instance
func (r *repository) Store(events <-chan *core.GamblerEvent) error {
	for {
		event, ok := <-events
		if !ok {
			return errors.New("Failed reading event")
		}
		log.Printf("Saving `%s` for `%s` on the `coins` Index", event.ID, event.Date)
		ctx := context.Background()
		data, err := json.Marshal(event)
		if err != nil {
			log.Fatal("Failed marshalling struct into JSON")
			return err
		}

		req := esapi.IndexRequest{
			Index:      "coins",
			DocumentID: fmt.Sprintf("%s_%s", event.ID, event.Date),
			Body:       bytes.NewReader(data),
			Refresh:    "true",
		}

		res, err := req.Do(ctx, r.client)
		if err != nil || res.IsError() {
			log.Fatalf("Failed indexing the data for: %v", event)
			return err
		}
		defer res.Body.Close()

		log.Print("Successfully inserted the event into the `coins` index")
	}
}
