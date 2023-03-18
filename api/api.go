package api

import (
	"deck_of_cards/card"
	"deck_of_cards/deck"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

var DeckStore *deck.Store

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/decks", createDeckHandler)
	r.GET("/decks/:deck_id", openDeckHandler)

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

	jsonResponse := DeckResponse{
		DeckID:    createdDeck.ID,
		Shuffled:  shuffled,
		Remaining: createdDeck.Remaining,
	}
	c.JSON(http.StatusOK, jsonResponse)
}

type DeckResponse struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
}

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

type OpenDeckResponse struct {
	DeckID    uuid.UUID   `json:"deck_id"`
	Shuffled  bool        `json:"shuffled"`
	Remaining int         `json:"remaining"`
	Cards     []card.Card `json:"cards"`
}
