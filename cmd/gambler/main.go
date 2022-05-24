package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nanchano/gambler/internal/core"
	"github.com/nanchano/gambler/pkg/pipeline/coingecko"
	"github.com/nanchano/gambler/pkg/repository/elastic"
)

func main() {
	coin := "ethereum"
	start := "20-04-2022"
	end := "20-05-2022"
	dates := createDateRange(start, end)
	pipeline := coingecko.NewPipeline(coin)
	repo := elastic.NewElasticRepository()
	service := core.NewGamblerService(pipeline, repo)

	getData(service, dates...)

	event, err := service.Find(coin, "25-04-2022")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%v\n", event)

}

func createDateRange(start, end string) []string {
	var dates []string
	startDate, _ := time.Parse("02-01-2006", start)
	endDate, _ := time.Parse("02-01-2006", end)
	for d := startDate; d.After(endDate) == false; d = d.AddDate(0, 0, 1) {
		ds := d.Format("02-01-2006")
		dates = append(dates, ds)
	}
	return dates
}

func getData(service core.GamblerService, dates ...string) {
	responses := service.Extract(dates...)
	events := service.Process(responses)
	service.Store(events)
}
