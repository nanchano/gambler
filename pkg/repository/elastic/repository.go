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

type elasticRepository struct {
	client *elastic.Client
}

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

// NewElasticRepository creates a new elasticRepository that implements core.GamblerRepository
func NewElasticRepository() core.GamblerRepository {
	client, err := newElasticClient()
	if err != nil {
		log.Fatalf("Error initializing the ElasticSearch client: %s", err)
	}

	return &elasticRepository{
		client: client,
	}
}

func (er *elasticRepository) Find(coin, date string) (*core.GamblerEvent, error) {
	var ge core.GamblerEvent
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

	res, err := er.client.Search(
		er.client.Search.WithContext(ctx),
		er.client.Search.WithIndex("coins"),
		er.client.Search.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error getting the response for %s - %s: %s", coin, date, err)
		return nil, err
	}
	defer res.Body.Close()

	var r elasticSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return nil, err
	}

	if len(r.Hits.Hits) == 1 {
		if err := json.Unmarshal(r.Hits.Hits[0].Source, &ge); err != nil {
			return nil, err
		}
		return &ge, nil
	} else {
		// TBD: implement multiple results
		return nil, errors.New("More than one result for the given filters")
	}

}

func (er *elasticRepository) Store(event *core.GamblerEvent) error {
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

	res, err := req.Do(ctx, er.client)
	if err != nil || res.IsError() {
		log.Fatalf("Failed indexing the data for: %v", event)
		return err
	}
	defer res.Body.Close()

	log.Printf("Successfully inserted %v into the `coins` index", event)
	return nil
}
