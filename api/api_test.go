package api_test

import (
	"deck_of_cards/api"
	"deck_of_cards/deck"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateDeckHandler(t *testing.T) {
	// Set up the Gin router and test server
	gin.SetMode(gin.TestMode) // Lightweight mode for testing.
	router := api.SetupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Perform a POST request to the /decks endpoint.
	resp, err := http.Post(server.URL+"/decks", "application/json", nil)
	assert.NoError(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// Decode the response body into a Deck.
	var createdDeck deck.Deck
	err = json.NewDecoder(resp.Body).Decode(&createdDeck)
	assert.NoError(t, err)

	// Assert that the created deck has the expected number of cards.
	assert.Equal(t, createdDeck.Remaining, 52)
	// Created deck has an ID.
	assert.NotNil(t, createdDeck.ID)
	// Created deck has 52 cards.
	assert.Equal(t, len(createdDeck.Cards), 52)
}
