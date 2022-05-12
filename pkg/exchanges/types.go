package exchanges

type Response map[string]interface{}

type Extractor func() <-chan []byte

type Processor func(<-chan []byte) <-chan *Response

type Consumer func(<-chan *Response)
