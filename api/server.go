// Package api provides the HTTP API for working with decks of playing cards.
// It uses the Gin web framework to handle HTTP requests and the `deck` and `card`
// packages to create and manage decks of cards. The package exposes endpoints
// for creating decks, opening decks, and drawing cards from decks.
package api

import (
	"deck_of_cards/deck"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *deck.Store
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{store: deck.NewStore()}
	router := gin.Default()

	router.POST("/deck/new", server.createDeckHandler)
	router.GET("/deck/:deck_id", server.openDeckHandler)
	router.POST("/deck/:deck_id/draw", server.drawCardHandler)

	server.router = router

	return server
}

func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
