package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type StockPrice struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}

func main() {
	http.HandleFunc("/ticker", stockTickerHandler)
	http.ListenAndServe(":8080", nil)
}

func stockTickerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive") // FIX 1: Changed "no-cache" to "keep-alive"
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return // FIX 2: Added return to stop execution on error
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()
	fmt.Println("client is connected to stock server")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Client disconnected from stock stream")
			return

		case <-ticker.C:
			data := StockPrice{
				Ticker: "AAPL",
				Price:  150.0 + (rand.Float64() * 10) - 5,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				continue
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if err != nil {
				return
			}
			flusher.Flush() // FIX 3: Moved inside the ticker case to flush only when writing data
		}
	}
}
