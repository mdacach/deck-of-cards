package api

import (
	"deck_of_cards/deck"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

// createDeckHandler is a Gin route handler for creating a new deck of cards.
// It accepts optional query parameters "cards" and "shuffled" to create a custom deck and shuffle it, respectively.
//
// Example query parameters for creating a partial deck and shuffling it:
// /decks?cards=AS,KD,QH,2C,3S&shuffled=true
//
// The deck information is returned as JSON.
func (server *Server) createDeckHandler(c *gin.Context) {
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

	err = server.store.Add(&createdDeck)
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
