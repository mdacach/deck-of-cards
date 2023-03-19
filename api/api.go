// Package api provides the HTTP API for working with decks of playing cards.
// It uses the Gin web framework to handle HTTP requests and the deck and card
// packages to create and manage decks of cards. The package exposes endpoints
// for creating decks, opening decks, and drawing cards from decks.
package api

import (
	"deck_of_cards/card"
	"deck_of_cards/deck"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

var DeckStore *deck.Store

// SetupRouter initializes the Gin router with the necessary routes for the card deck API.
// It returns a pointer to the router, which can be used to start the HTTP server.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/deck/new", createDeckHandler)
	r.GET("/deck/:deck_id", openDeckHandler)
	r.GET("/deck/:deck_id/draw", drawCardHandler)

	return r
}

// createDeckHandler is a Gin route handler for creating a new deck of cards.
// It accepts optional query parameters "cards" and "shuffled" to create a custom deck and shuffle it, respectively.
//
// Example query parameters for creating a partial deck and shuffling it:
// /decks?cards=AS,KD,QH,2C,3S&shuffled=true
//
// The deck information is returned as JSON.
func createDeckHandler(c *gin.Context) {
	queryCards, exists := c.GetQuery("cards")
	var createdDeck deck.Deck
	var err error
	if exists {
		cardCodes := strings.Split(queryCards, ",")
		createdDeck, err = deck.NewPartialDeck(cardCodes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		createdDeck = deck.NewStandardDeck()
	}

	// TODO: Some better, strongly typed way of doing this?
	shuffledStr := c.DefaultQuery("shuffled", "false")
	shuffled := false
	if shuffledStr == "true" {
		shuffled = true
		createdDeck.Shuffle()
	}

	err = DeckStore.Add(&createdDeck)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	jsonResponse := CreateDeckResponse{
		DeckID:    createdDeck.ID,
		Shuffled:  shuffled,
		Remaining: createdDeck.Remaining,
	}
	c.JSON(http.StatusOK, jsonResponse)
}

// CreateDeckResponse is a struct that represents the JSON response for the createDeckHandler.
type CreateDeckResponse struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
}

// openDeckHandler is a Gin route handler for retrieving an existing deck by its ID.
// The deck ID is provided as a URL parameter. If the deck is found, the deck information is returned as JSON.
func openDeckHandler(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("deck_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck ID is not valid."})
	}

	deckRetrieved, notFound := DeckStore.Get(deckID)
	if notFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck not found. Are you sure deck_id is correct?"})
	}

	jsonResponse := OpenDeckResponse{
		DeckID:    deckRetrieved.ID,
		Shuffled:  deckRetrieved.Shuffled,
		Remaining: deckRetrieved.Remaining,
		Cards:     deckRetrieved.Cards,
	}
	c.JSON(http.StatusOK, jsonResponse)
}

// OpenDeckResponse is a struct that represents the JSON response for the openDeckHandler.
type OpenDeckResponse struct {
	DeckID    uuid.UUID   `json:"deck_id"`
	Shuffled  bool        `json:"shuffled"`
	Remaining int         `json:"remaining"`
	Cards     []card.Card `json:"cards"`
}

// drawCardHandler is a Gin route handler for drawing a specified number of cards from an existing deck.
// The deck ID and card count are provided as URL parameters. If the deck is found and the draw is successful,
// the drawn cards are returned as JSON.
func drawCardHandler(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("deck_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck ID is not valid."})
		return
	}

	countStr, exists := c.GetQuery("count")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count parameter must be provided."})
		return
	}
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count parameter must be a positive integer"})
		return
	}

	deckRetrieved, notFound := DeckStore.Get(deckID)
	if notFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck not found. Are you sure deck_id is correct?"})
		return
	}

	drawnCards, err := deckRetrieved.Draw(count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	jsonResponse := DrawCardsResponse{
		Cards: drawnCards,
	}
	c.JSON(http.StatusOK, jsonResponse)
}

// DrawCardsResponse is a struct that represents the JSON response for the drawCardHandler.
type DrawCardsResponse struct {
	Cards []card.Card `json:"cards"`
}
