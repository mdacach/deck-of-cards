package api_test

import (
	"bytes"
	"deck_of_cards/api"
	"deck_of_cards/card"
	"deck_of_cards/deck"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCreatePartialDeckEndpoint(t *testing.T) {
	router := api.SetupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	cards := "AS,KD,AC,2C,KH"
	resp, err := http.Post(ts.URL+"/decks?cards="+cards, "", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var createdDeck deck.Deck
	err = json.NewDecoder(resp.Body).Decode(&createdDeck)
	require.NoError(t, err)

	expectedCards := []card.Card{
		{Rank: "A", Suit: "S"},
		{Rank: "K", Suit: "D"},
		{Rank: "A", Suit: "C"},
		{Rank: "2", Suit: "C"},
		{Rank: "K", Suit: "H"},
	}

	assert.Equal(t, len(expectedCards), len(createdDeck.Cards))
	assert.Equal(t, len(expectedCards), createdDeck.Remaining)

	for i, c := range expectedCards {
		assert.Equal(t, c, createdDeck.Cards[i])
	}
}

func TestCreateDeckHandlerInvalidRequests(t *testing.T) {
	router := api.SetupRouter()

	testCases := []struct {
		name       string
		cardsParam string
	}{
		{
			name:       "cards param with no cards",
			cardsParam: "",
		},
		{
			name:       "cards with invalid card code",
			cardsParam: "INVALID_CARD",
		},
		{
			name:       "cards with repeated codes",
			cardsParam: "AS,AS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			req, _ := http.NewRequest("POST", "/decks?cards="+tc.cardsParam, body)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestCreateStandardDeckShuffled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := api.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/decks?shuffled=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp api.DeckResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// From DeckResponse we do not have access to the cards directly,
	// but let's assert that the shuffled requirement was set.
	assert.NotEmpty(t, resp.DeckID)
	assert.True(t, resp.Shuffled)
	assert.Equal(t, 52, resp.Remaining)
}
