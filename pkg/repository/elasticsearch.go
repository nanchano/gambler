package repository

import (
	"io"
	"io/ioutil"
	"log"

	elastic "github.com/elastic/go-elasticsearch/v8"
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

func NewElasticRepository() core.GamblerRepository {
	client, err := newElasticClient()
	if err != nil {
		log.Fatalf("Error initializing the ElasticSearch client: %s", err)
	}

	return elasticRepository{
		client: client,
	}
}

func (er *elasticRepository) Store(ge *core.GamblerEvent) error {

}
