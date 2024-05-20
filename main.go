package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LTPResponse struct {
	LTP []struct {
		Pair   string  `json:"pair"`
		Amount float64 `json:"amount"`
	} `json:"ltp"`
}

// Makes a request to the 3rd-party API and returns the response.
func getLTP() (map[string]interface{}, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=BTCCHF,BTCEUR,BTCUSD"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Processes the response from the 3rd-party API and creates the LTPResponse structure.
func processResponse(response map[string]interface{}) LTPResponse {
	result := response["result"].(map[string]interface{})

	ltp := make([]struct {
		Pair   string  `json:"pair"`
		Amount float64 `json:"amount"`
	}, 0, len(result))

	for pair, data := range result {
		price := data.(map[string]interface{})["c"].([]interface{})[0].(string)
		log.Print(price)
		ltp = append(ltp, struct {
			Pair   string  `json:"pair"`
			Amount float64 `json:"amount"`
		}{
			Pair:   formatPair(pair),
			Amount: parseFloat(price),
		})
	}

	return LTPResponse{LTP: ltp}
}

// Formats the pair string as "BTC/CHF", "BTC/EUR", or "BTC/USD".
func formatPair(pair string) string {
	quote := pair[len(pair)-3:]
	return fmt.Sprintf("BTC/%s", quote)
}

// Parses the price string into a float64.
func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func main() {
	r := gin.Default()

	r.GET("/api/v1/ltp", func(c *gin.Context) {
		// Call the 3rd-party API and get the response
		response, err := getLTP()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Process the response and create the LTPResponse
		ltpResponse := processResponse(response)

		c.JSON(http.StatusOK, ltpResponse)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	r.Run(":" + port)
}
