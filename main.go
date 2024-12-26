package main

import (
	"Squeeze_Alert/indicator"
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	closePrices := []float64{10, 12, 15, 14, 13, 11, 10, 9, 10, 11, 13, 14, 15, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 8, 9}
	highPrices := []float64{11, 13, 16, 15, 14, 12, 11, 10, 11, 12, 14, 15, 16, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 9, 10}
	lowPrices := []float64{9, 11, 14, 13, 12, 10, 9, 8, 9, 10, 12, 13, 14, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 7, 8}

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		ind := indicator.NewSqueezeIndicator(20, 20, 2.0, 1.5, true, 0)
		values, sqzOn := ind.Calculate(closePrices, highPrices, lowPrices)

		data := make([]map[string]interface{}, len(values))
		for i := range values {
			data[i] = map[string]interface{}{
				"time":  i,
				"close": closePrices[i],
				"value": values[i],
				"sqzOn": sqzOn[i],
			}
		}

		json.NewEncoder(w).Encode(data)
	})

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fs)

	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
