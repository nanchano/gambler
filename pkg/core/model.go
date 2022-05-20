package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// GamblerEvent represents a normalized response from any crypto API
type GamblerEvent struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Symbol    string      `json:"symbol"`
	Date      string      `json:"date"`
	Price     float64     `json:"price"`
	MarketCap float64     `json:"market_cap"`
	Volume    float64     `json:"volume"`
	Extra     interface{} `json:"extra"`
}

func (gr *GamblerEvent) ToJSON() error {
	outputDir := fmt.Sprintf("outputs/coingecko/%s", gr.ID)
	outputPath := fmt.Sprintf("%v/%v.json", outputDir, gr.Date)
	log.Printf("Saving to %v", outputPath)

	file, _ := json.MarshalIndent(gr, "", "\t")
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalln("Failed creating the directories")
		return err
	}
	err = ioutil.WriteFile(outputPath, file, 0644)
	if err != nil {
		log.Fatalln("Failed riting the JSON")
		return err
	}

	return nil
}
