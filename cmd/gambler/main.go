package main

import (
	"fmt"
	"time"

	"github.com/nanchano/gambler/pkg/exchanges"
)

func main() {
	id := "ethereum"
	// date := "20-04-2022"
	// h := exchanges.NewCoingeckoHandler(id, date)
	// h.ResponseConsumer()(h.ResponseProcessor()(h.ResponseGenerator()()))
	// start := time.Now()
	// end := start.AddDate(0, 1, 0)
	start, _ := time.Parse("02-01-2006", "01-04-2022")
	end, _ := time.Parse("02-01-2006", "01-05-2022")
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		ds := d.Format("02-01-2006")
		fmt.Println(ds)
		h := exchanges.NewCoingeckoHandler(id, ds)
		h.ResponseConsumer()(h.ResponseProcessor()(h.ResponseGenerator()()))
	}
}
