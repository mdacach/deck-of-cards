package api

import (
	"deck_of_cards/card"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenDeck(t *testing.T) {
	router := setup()

	// 1. Create the deck through the Create endpoint. It will be stored (somewhere).
	// Create a new standard deck using the Create Deck endpoint.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/deck/new", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusOK)

	var createResponse CreateDeckResponse
	err := json.NewDecoder(w.Body).Decode(&createResponse)
	require.NoError(t, err)

	// Keep track of the deck's ID.
	deckID := createResponse.DeckID

	// 2. Open the (same) deck through Open endpoint.
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/deck/%s", deckID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	var openResponse OpenDeckResponse
	err = json.NewDecoder(w.Body).Decode(&openResponse)
	require.NoError(t, err)

	assert.Equal(t, deckID, openResponse.DeckID, "Deck ID does not change after it is created.")
	assert.False(t, openResponse.Shuffled, "Deck Shuffled does not change after it is created.")
	assert.Equal(t, 52, openResponse.Remaining, "If we do not Draw from the deck, all cards still remain.")
}

func TestOpenPartialDeck(t *testing.T) {
	router := setup()

	cardCodes := "AS,KD,AC,2C,KH"
	expectedCards := []card.Card{
		{Rank: "A", Suit: "S"},
		{Rank: "K", Suit: "D"},
		{Rank: "A", Suit: "C"},
		{Rank: "2", Suit: "C"},
		{Rank: "K", Suit: "H"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/deck/new?cards="+cardCodes, nil)
	router.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusOK)

	var createResponse CreateDeckResponse
	err := json.NewDecoder(w.Body).Decode(&createResponse)
	require.NoError(t, err)

	// Keep track of the deck's ID.
	deckID := createResponse.DeckID

	// 2. Open the (same) deck through Open endpoint.
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/deck/%s", deckID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	var openResponse OpenDeckResponse
	err = json.NewDecoder(w.Body).Decode(&openResponse)
	require.NoError(t, err)

	assert.Equal(t, deckID, openResponse.DeckID, "Deck ID does not change after it is created.")
	assert.False(t, openResponse.Shuffled, "Deck Shuffled does not change after it is created.")
	assert.Equal(t, len(expectedCards), openResponse.Remaining, "If we do not Draw from the deck, all cards still remain.")

	// The cards are in the correct order (the order we specified in the request).
	for i, c := range expectedCards {
		assert.Equal(t, c, openResponse.Cards[i])
	}
}
