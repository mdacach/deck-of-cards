package api

import (
	"deck-of-cards/card"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateDeckHandler(t *testing.T) {
	router := setup()

	// Perform a POST request to the /decks endpoint.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/deck/new", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	// Decode the response body.
	var createdDeck CreateDeckResponse
	err := json.NewDecoder(w.Body).Decode(&createdDeck)
	assert.NoError(t, err)

	// Assert that the created deck has the expected number of cards.
	assert.Equal(t, createdDeck.Remaining, 52)
	// Created deck has an ID.
	assert.NotNil(t, createdDeck.DeckID)
}

func TestCreatePartialDeckEndpoint(t *testing.T) {
	router := setup()

	cards := "AS,KD,AC,2C,KH"
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/deck/new?cards="+cards, nil)
	router.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusOK)

	var createdDeck CreateDeckResponse
	err := json.NewDecoder(w.Body).Decode(&createdDeck)
	require.NoError(t, err)

	expectedCards := []card.Card{
		{Rank: "A", Suit: "S"},
		{Rank: "K", Suit: "D"},
		{Rank: "A", Suit: "C"},
		{Rank: "2", Suit: "C"},
		{Rank: "K", Suit: "H"},
	}

	assert.Equal(t, len(expectedCards), createdDeck.Remaining)
}

func TestCreateDeckHandlerInvalidRequests(t *testing.T) {
	router := setup()

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
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/deck/new?cards="+tc.cardsParam, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, http.StatusBadRequest)
		})
	}
}

func TestCreateStandardDeckShuffled(t *testing.T) {
	router := setup()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/deck/new?shuffled=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	var resp CreateDeckResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)

	// From CreateDeckResponse we do not have access to the cards directly,
	// but let's assert that the shuffled requirement was set.
	assert.NotEmpty(t, resp.DeckID)
	assert.True(t, resp.Shuffled)
	assert.Equal(t, 52, resp.Remaining)
}
