package api

import (
	"deck-of-cards/card"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDrawCardHandler(t *testing.T) {
	router := setup()

	validID := createTestDeck(router, "").String()

	testCases := []struct {
		name         string
		deckID       string
		count        string
		expectedCode int
	}{
		{
			name:         "invalid deck ID",
			deckID:       "invalid-deck-id",
			count:        "5",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "no count parameter provided",
			deckID:       validID,
			count:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid count parameter",
			deckID:       validID,
			count:        "invalid-count",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "non-positive count parameter",
			deckID:       validID,
			count:        "-5",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "deck not found",
			deckID:       uuid.NewString(),
			count:        "5",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target := fmt.Sprintf("/deck/%s/draw?count=%s", tc.deckID, tc.count)
			req := httptest.NewRequest(http.MethodPost, target, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedCode, "Expected status code to match")
		})
	}
}

func TestDrawPartialDeck(t *testing.T) {
	router := setup()

	cardCodes := "QH,4D,AC,2C,KH"
	expectedCards := []card.Card{
		{Rank: "Q", Suit: "H"},
		{Rank: "4", Suit: "D"},
		{Rank: "A", Suit: "C"},
		{Rank: "2", Suit: "C"},
		{Rank: "K", Suit: "H"},
	}

	deckID := createTestDeck(router, "?cards="+cardCodes)

	// Draw the first card: it should be the Queen of Hearts (QH).
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/deck/%s/draw?count=1", deckID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var drawResponse DrawCardsResponse
	err := json.NewDecoder(w.Body).Decode(&drawResponse)
	assert.NoError(t, err)
	drawnCards := drawResponse.Cards
	assert.Equal(t, len(drawnCards), 1, "One card is drawn.")
	assert.Equal(t, drawnCards[0], expectedCards[0], "Card drawn is the first in the deck.")

	// Draw a new card: it should be the 4 of Diamonds (4D).
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/deck/%s/draw?count=1", deckID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.NewDecoder(w.Body).Decode(&drawResponse)
	assert.NoError(t, err)
	drawnCards = drawResponse.Cards
	assert.Equal(t, len(drawnCards), 1, "One card is drawn.")
	assert.Equal(t, drawnCards[0], expectedCards[1], "Card drawn is the (currently) first in the deck.")

	// Draw the three last cards.
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/deck/%s/draw?count=3", deckID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.NewDecoder(w.Body).Decode(&drawResponse)
	assert.NoError(t, err)

	drawnCards = drawResponse.Cards
	assert.Equal(t, len(drawnCards), 3, "Three cards are drawn.")
	assert.Equal(t, drawnCards[0], expectedCards[2], "Card drawn is the (currently) first in the deck.")
	assert.Equal(t, drawnCards[1], expectedCards[3], "Next card drawn is the (currently) first in the deck.")
	assert.Equal(t, drawnCards[2], expectedCards[4], "Next card drawn is the (currently) first in the deck.")
}
