package main

import (
	"fmt"
	"time"

	"github.com/nanchano/gambler/pkg/exchanges"
	"github.com/nanchano/gambler/pkg/repository/elastic"
)

func main() {
	// id := "ethereum"
	// dates := createDateRange("01-04-2022", "01-05-2022")
	// h := exchanges.NewCoingeckoPipeline(id)
	// run(h, dates...)

	// ge := core.GamblerEvent{
	// 	ID:        "ethereum",
	// 	Name:      "Ethereum",
	// 	Symbol:    "eth",
	// 	Date:      "20-04-2023",
	// 	Price:     2000.01,
	// 	MarketCap: 1.66,
	// 	Volume:    2.66,
	// 	Extra:     "BCE",
	// }

	er := elastic.NewElasticRepository()
	// er.Store(&ge)
	ge, err := er.Find("ethereum", "20-04-2022")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("\n%v", ge)
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

func run(p *exchanges.CoingeckoPipeline, dates ...string) {
	in := p.Extract(dates...)
	out := p.Process(in)
	p.Save(out)
}
