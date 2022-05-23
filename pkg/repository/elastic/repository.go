package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/nanchano/gambler/internal/core"
)

type ElasticRepository struct {
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
// func NewElasticRepository() core.GamblerRepository {
func NewElasticRepository() ElasticRepository {
	client, err := newElasticClient()
	if err != nil {
		log.Fatalf("Error initializing the ElasticSearch client: %s", err)
	}

	return ElasticRepository{
		client: client,
	}
}

func (er *ElasticRepository) Find(coin, date string) (*core.GamblerEvent, error) {
	var ge core.GamblerEvent

	// "query": {
	// 	"query_string": {
	// 		"query": "(id:\"ethereum\") AND (date:\"20-04-2022\")"
	// 	}
	// }

	ctx := context.Background()
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": fmt.Sprintf("(id:\"ethereum\") AND (date:\"20-04-2022\")"),
			},
		},
	}
	_ = json.NewEncoder(&buf).Encode(query)

	res, err := er.client.Search(
		er.client.Search.WithContext(ctx),
		er.client.Search.WithIndex("coins"),
		er.client.Search.WithBody(&buf),
		// er.client.Search.WithTrackTotalHits(true),
		er.client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting the response for %s - %s: %s", coin, date, err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatal("Error parsing the response body: %s", err)
		return nil, err
	}

	// Print the ID and document source for each hit.
	for i, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		if i == 0 {
			continue
		}

		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])

		if err := json.NewDecoder(hit.(map[string]interface{})["_source"]).Decode(&ge); err != nil {
			log.Fatal("Error parsing the response body: %s", err)
			return nil, err
		}

	}

	log.Printf("%v", ge)

	return &ge, nil

}

func (er *ElasticRepository) Store(ge *core.GamblerEvent) error {
	log.Printf("Saving `%s` for `%s` on the `coins` Index", ge.ID, ge.Date)
	ctx := context.Background()

	data, err := json.Marshal(ge)
	if err != nil {
		log.Fatal("Failed marshalling struct into JSON")
		return err
	}

	req := esapi.IndexRequest{
		Index:      "coins",
		DocumentID: fmt.Sprintf("%s_%s", ge.ID, ge.Date),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, er.client)
	if err != nil || res.IsError() {
		log.Fatal("Failed indexing the data for: %v", ge)
		return err
	}
	defer res.Body.Close()

	log.Printf("Successfully inserted %v into the `coins` index", ge)
	return nil
}
