package elastic

import "encoding/json"

type innerHit struct {
	ID         string          `json:"_id"`
	Source     json.RawMessage `json:"_source"`
	Highlights json.RawMessage `json:"highlight"`
	Sort       []interface{}   `json:"sort"`
}

type total struct {
	Value int
}

type outerHits struct {
	Total total
	Hits  []innerHit
}

type elasticSearchResponse struct {
	Took int
	Hits outerHits
}
