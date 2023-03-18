package api

import (
	"deck_of_cards/deck"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/decks", createDeckHandler)

	return r
}

func createDeckHandler(c *gin.Context) {
	d := deck.NewStandardDeck()
	c.JSON(http.StatusOK, d)
}
