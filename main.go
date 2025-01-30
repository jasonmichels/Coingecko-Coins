package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Coin struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	MarketCapRank int    `json:"market_cap_rank"`
}

func main() {
	var allCoins []Coin
	maxPages := 16

	for page := 1; page <= maxPages; page++ {
		url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd" +
			"&order=market_cap_desc&per_page=250&page=" + strconv.Itoa(page)

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			// Read error details
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			panic(fmt.Sprintf("Non-200 response: %d\nBody: %s", resp.StatusCode, body))
		}

		var pageCoins []Coin
		if err := json.NewDecoder(resp.Body).Decode(&pageCoins); err != nil {
			resp.Body.Close()
			panic(err)
		}
		resp.Body.Close()

		// Filter for MarketCapRank <= 4000
		for _, c := range pageCoins {
			if c.MarketCapRank > 0 && c.MarketCapRank <= 4000 {
				allCoins = append(allCoins, c)
			}
		}

		// If fewer than 250 results were returned, we’re likely done
		if len(pageCoins) < 250 {
			break
		}

		// 20s delay between pages to avoid 429s
		time.Sleep(20 * time.Second)
	}

	// Write results to JSON file
	f, err := os.Create("coins.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(allCoins); err != nil {
		panic(err)
	}

	fmt.Printf("Saved %d coins to coins.json\n", len(allCoins))
}
