// Package main provides the entry point for the deck_of_cards application, which
// is a REST API for managing and manipulating decks of playing cards.
//
// The application uses the Gin framework for handling HTTP requests, and it
// exposes the following endpoints:
// - POST /decks: Create a new deck, either a standard deck or a custom one with specified cards, and shuffle it if needed
// - GET /decks/:deck_id: Retrieve the information of an existing deck
// - GET /decks/:deck_id/draw: Draw a specified number of cards from an existing deck
//
// The API is served on port 8080 by default.
package main

import (
	"deck_of_cards/api"
	"deck_of_cards/deck"
	"fmt"
)

func main() {
	// TODO: Change this to (probably) Dependency Injection.
	api.DeckStore = deck.NewStore()
	r := api.SetupRouter()
	r.Run(":8080")
	router := api.SetupRouter()

	// TODO: Get port to run from flag/env variable
	err := router.Run(":8080")
	if err != nil {
		// TODO: Better error handling.
		fmt.Println("Could not start Gin server. Maybe port :8080 is already being used?")
		return
	}
}
