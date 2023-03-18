package api_test

import (
	"bytes"
	"deck_of_cards/api"
	"deck_of_cards/card"
	"deck_of_cards/deck"
	"encoding/json"
	"fmt"
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

func TestOpenDeck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := api.SetupRouter()
	// TODO: Improve this. It's not nice to need to set this global variable every time.
	//       Will probably remove the global variable, but if not, at least create a setup function for tests.
	api.DeckStore = deck.NewStore()

	// 1. Create the deck through the Create endpoint. It will be stored (somewhere).
	// Create a new standard deck using the Create Deck endpoint.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/decks", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var createResponse api.DeckResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)

	// Keep track of the deck's ID.
	deckID := createResponse.DeckID

	// 2. Open the (same) deck through Open endpoint.
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/decks/%s", deckID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var openResponse api.OpenDeckResponse
	err = json.Unmarshal(w.Body.Bytes(), &openResponse)
	require.NoError(t, err)

	assert.Equal(t, deckID, openResponse.DeckID, "Deck ID does not change after it is created.")
	assert.False(t, openResponse.Shuffled, "Deck Shuffled does not change after it is created.")
	assert.Equal(t, 52, openResponse.Remaining, "If we do not Draw from the deck, all cards still remain.")
}

func TestOpenPartialDeck(t *testing.T) {
	router := api.SetupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	api.DeckStore = deck.NewStore()

	cardCodes := "AS,KD,AC,2C,KH"
	expectedCards := []card.Card{
		{Rank: "A", Suit: "S"},
		{Rank: "K", Suit: "D"},
		{Rank: "A", Suit: "C"},
		{Rank: "2", Suit: "C"},
		{Rank: "K", Suit: "H"},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/decks?cards="+cardCodes, nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var createResponse api.DeckResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)

	// Keep track of the deck's ID.
	deckID := createResponse.DeckID

	// 2. Open the (same) deck through Open endpoint.
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/decks/%s", deckID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var openResponse api.OpenDeckResponse
	err = json.Unmarshal(w.Body.Bytes(), &openResponse)
	require.NoError(t, err)

	assert.Equal(t, deckID, openResponse.DeckID, "Deck ID does not change after it is created.")
	assert.False(t, openResponse.Shuffled, "Deck Shuffled does not change after it is created.")
	assert.Equal(t, len(expectedCards), openResponse.Remaining, "If we do not Draw from the deck, all cards still remain.")

	// The cards are in the correct order (the order we specified in the request).
	for i, c := range expectedCards {
		assert.Equal(t, c, openResponse.Cards[i])
	}
}
