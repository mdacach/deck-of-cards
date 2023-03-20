package api

import (
	"deck-of-cards/card"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

// drawCardHandler is a Gin route handler for drawing a specified number of cards from an existing deck.
// The deck ID and card count are provided as URL parameters. If the deck is found and the draw is successful,
// the drawn cards are returned as JSON.
func (server *Server) drawCardHandler(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("deck_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deck ID is not valid."})
		return
	}

	countStr, exists := c.GetQuery("count")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "count parameter must be provided."})
		return
	}
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "count parameter must be a positive integer"})
		return
	}

	deckRetrieved, notFound := server.store.Get(deckID)
	if notFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deck not found. Are you sure deck_id is correct?"})
		return
	}

	drawnCards, err := deckRetrieved.Draw(count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
