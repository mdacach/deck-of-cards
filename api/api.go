package api

import (
	"deck_of_cards/deck"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/decks", createDeckHandler)

	return r
}

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

	jsonResponse := DeckResponse{
		DeckID:    createdDeck.ID,
		Shuffled:  false,
		Remaining: createdDeck.Remaining,
	}
	c.JSON(http.StatusOK, jsonResponse)
}

type DeckResponse struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
}
