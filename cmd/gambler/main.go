package main

import (
	"fmt"

	"github.com/nanchano/gambler/pkg/exchanges"
)

func main() {
	id := "ethereum"
	date := "20-04-2022"
	h := exchanges.NewCoingeckoHandler(id, date)
	b := h.Extract()
	t := h.Transform(b)
	fmt.Println(t)
}
