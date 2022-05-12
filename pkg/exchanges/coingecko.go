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

func (c *CoingeckoHandler) prepareURL() string {
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

func (c *CoingeckoHandler) ResponseGenerator() Extractor {
	fmt.Println("Outer call")
	return func() <-chan []byte {
		fmt.Println("Inner call")
		out := make(chan []byte)
		go func() {
			fmt.Println("Go call")
			defer close(out)
			url := c.prepareURL()
			resp, err := http.Get(url)

			if err != nil {
				log.Fatalln(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("ASD")
			fmt.Println(string(body))
			out <- body
			fmt.Println("ASD")
			fmt.Println(string(body))
		}()
		return out
	}
}

func (c *CoingeckoHandler) ResponseProcessor() Processor {
	return func(in <-chan []byte) <-chan *Response {
		out := make(chan *Response)
		go func() {
			defer close(out)
			for resp := range in {
				var r Response
				json.Unmarshal(resp, &r)
				out <- &r
			}
		}()
		return out
	}
}

func (c *CoingeckoHandler) ResponseConsumer() Consumer {
	return func(in <-chan *Response) {
		for {
			i, ok := <-in
			if ok {
				fmt.Printf("%v\n", i)
			} else {
				return
			}
		}
	}
}
