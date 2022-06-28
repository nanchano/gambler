package elastic

import "encoding/json"

// innerHit corresponds to the innerHit key on the ElasticSearch JSON ersponse
type innerHit struct {
	ID         string          `json:"_id"`
	Source     json.RawMessage `json:"_source"`
	Highlights json.RawMessage `json:"highlight"`
	Sort       []interface{}   `json:"sort"`
}

// total corresponds to the total key on the ElasticSearch JSON response
type total struct {
	Value int
}

// outerHits corresponds to the outerHits key on the ElasticSearch JSON response
type outerHits struct {
	Total total
	Hits  []innerHit
}

// elasticSearchResponse represents a raw response to a query on an ElasticSearch instance
type elasticSearchResponse struct {
	Took int
	Hits outerHits
}
