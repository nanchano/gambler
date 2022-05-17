package main

import (
	"time"

	"github.com/nanchano/gambler/pkg/exchanges"
)

func main() {
	id := "ethereum"
	dates := createDateRange("01-04-2022", "01-05-2022")
	h := exchanges.NewCoingeckoHandler(id)
	run(h, dates...)
}

func createDateRange(start, end string) []string {
	var dates []string
	start_date, _ := time.Parse("02-01-2006", start)
	end_date, _ := time.Parse("02-01-2006", end)
	for d := start_date; d.After(end_date) == false; d = d.AddDate(0, 0, 1) {
		ds := d.Format("02-01-2006")
		dates = append(dates, ds)
	}
	return dates
}

func run(h *exchanges.CoingeckoHandler, dates ...string) {
	in := h.Extract(dates...)
	out := h.Process(in)
	h.Save(out)
}
