package api

import (
	"deck_of_cards/card"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// openDeckHandler is a Gin route handler for retrieving an existing deck by its ID.
// The deck ID is provided as a URL parameter. If the deck is found, the deck information is returned as JSON.
func (server *Server) openDeckHandler(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("deck_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck ID is not valid."})
	}

	deckRetrieved, notFound := server.store.Get(deckID)
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
