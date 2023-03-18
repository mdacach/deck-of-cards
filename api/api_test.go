package api_test

import (
	"bytes"
	"deck_of_cards/api"
	"deck_of_cards/card"
	"deck_of_cards/deck"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func createTestDeck(router *gin.Engine, params string) uuid.UUID {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/decks"+params, nil)
	router.ServeHTTP(w, req)

	var createResponse api.DeckResponse
	// This should never fail.
	_ = json.Unmarshal(w.Body.Bytes(), &createResponse)

	deckID := createResponse.DeckID

	return deckID
}

func TestDrawCardHandler(t *testing.T) {
	router := api.SetupRouter()
	api.DeckStore = deck.NewStore()

	// Define test cases
	type testCase struct {
		name         string
		deckID       string
		count        string
		expectedCode int
	}

	validID := createTestDeck(router, "").String()

	testCases := []testCase{
		{
			name:         "Invalid deck ID",
			deckID:       "invalid-deck-id",
			count:        "5",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "No count parameter provided",
			deckID:       validID,
			count:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid count parameter",
			deckID:       validID,
			count:        "invalid-count",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Non-positive count parameter",
			deckID:       validID,
			count:        "-5",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Deck not found",
			deckID:       uuid.NewString(),
			count:        "5",
			expectedCode: http.StatusBadRequest,
		},
	}

	// Test the drawCardHandler
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the test data
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/draw?count=%s", tc.deckID, tc.count), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Execute the request and record the response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Use assert to check if the expected code matches the actual code
			assert.Equal(t, tc.expectedCode, rr.Code, "Expected status code to match")
		})
	}
}

func TestDrawPartialDeck(t *testing.T) {
	router := api.SetupRouter()
	api.DeckStore = deck.NewStore()

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
	req, _ := http.NewRequest("GET", fmt.Sprintf("/decks/%s/draw?count=1", deckID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var drawResponse api.DrawCardsResponse
	_ = json.Unmarshal(w.Body.Bytes(), &drawResponse)
	drawnCards := drawResponse.Cards
	assert.Equal(t, len(drawnCards), 1, "One card is drawn.")
	assert.Equal(t, drawnCards[0], expectedCards[0], "Card drawn is the first in the deck.")

	// Draw a new card: it should be the 4 of Diamonds (4D).
	req, _ = http.NewRequest("GET", fmt.Sprintf("/decks/%s/draw?count=1", deckID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_ = json.Unmarshal(w.Body.Bytes(), &drawResponse)
	drawnCards = drawResponse.Cards
	assert.Equal(t, len(drawnCards), 1, "One card is drawn.")
	assert.Equal(t, drawnCards[0], expectedCards[1], "Card drawn is the (currently) first in the deck.")

	// Draw the three last cards.
	req, _ = http.NewRequest("GET", fmt.Sprintf("/decks/%s/draw?count=3", deckID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_ = json.Unmarshal(w.Body.Bytes(), &drawResponse)
	drawnCards = drawResponse.Cards
	assert.Equal(t, len(drawnCards), 3, "Three cards are drawn.")
	assert.Equal(t, drawnCards[0], expectedCards[2], "Card drawn is the (currently) first in the deck.")
	assert.Equal(t, drawnCards[1], expectedCards[3], "Next card drawn is the (currently) first in the deck.")
	assert.Equal(t, drawnCards[2], expectedCards[4], "Next card drawn is the (currently) first in the deck.")
}
