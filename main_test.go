package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetLTP(t *testing.T) {
	response, err := getLTP()
	if err != nil {
		t.Errorf("Failed to get LTP: %v", err)
	}

	// Check if the response has the expected structure
	_, ok := response["result"]
	if !ok {
		t.Errorf("Invalid response structure")
	}

	// Check if the response contains the expected pairs
	expectedPairs := []string{"XBTCHF", "XXBTZEUR", "XXBTZUSD"}
	result := response["result"].(map[string]interface{})
	for _, pair := range expectedPairs {
		_, ok := result[pair]
		if !ok {
			t.Errorf("Missing pair %s in response", pair)
		}
	}
}

// Mock getLTP function that returns a predefined response
func mockGetLTP() (map[string]interface{}, error) {
	return map[string]interface{}{
		"result": map[string]interface{}{
			"XBTCHF": map[string]interface{}{
				"c": []interface{}{"1000.00"},
			},
			"XXBTZEUR": map[string]interface{}{
				"c": []interface{}{"900.00"},
			},
			"XXBTZUSD": map[string]interface{}{
				"c": []interface{}{"1100.00"},
			},
		},
	}, nil
}

// Function to get the index of a string in a slice
func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1 // Return -1 if the item is not found
}

func TestProcessResponse(t *testing.T) {
	// Mock response
	response, _ := mockGetLTP()

	ltpResponse := processResponse(response)

	// Check if the LTPResponse has the expected structure
	if len(ltpResponse.LTP) != len(response["result"].(map[string]interface{})) {
		t.Errorf("Invalid LTPResponse structure")
	}

	// Check if the LTPResponse contains the expected pairs and amounts
	expectedPairs := []string{"BTC/CHF", "BTC/EUR", "BTC/USD"}
	expectedAmounts := []float64{1000.00, 900.00, 1100.00}
	for _, ltp := range ltpResponse.LTP {
		idx := indexOf(expectedPairs, ltp.Pair)

		if idx == -1 {
			t.Errorf("Invalid pair in LTPResponse")
		}
		if ltp.Amount != expectedAmounts[idx] {
			t.Errorf("Invalid amount in LTPResponse")
		}
	}
}

func TestAPILTP(t *testing.T) {
	// Create a new Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Mount the API endpoint
	r.GET("/api/v1/ltp", func(c *gin.Context) {
		// Call the mock getLTP function and get the response
		response, err := mockGetLTP()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Process the response and create the LTPResponse
		ltpResponse := processResponse(response)

		c.JSON(http.StatusOK, ltpResponse)
	})

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/api/v1/ltp", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Create a new HTTP recorder to capture the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Read the response body
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the response body into a LTPResponse struct
	var ltpResponse LTPResponse
	err = json.Unmarshal(body, &ltpResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check if the LTPResponse has the expected structure
	if len(ltpResponse.LTP) == 0 {
		t.Errorf("Expected non-empty LTPResponse")
	}

	// Check if the LTPResponse contains the expected pairs and amounts
	expectedPairs := []string{"BTC/CHF", "BTC/EUR", "BTC/USD"}
	expectedAmounts := []float64{1000.00, 900.00, 1100.00}
	for _, ltp := range ltpResponse.LTP {
		idx := indexOf(expectedPairs, ltp.Pair)

		if idx == -1 {
			t.Errorf("Invalid pair in LTPResponse")
		}
		if ltp.Amount != expectedAmounts[idx] {
			t.Errorf("Expected amount %f, got %f", expectedAmounts[idx], ltp.Amount)
		}
	}

}
