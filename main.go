package main

import (
	"deck_of_cards/api"
	"deck_of_cards/deck"
)

func main() {
	api.DeckStore = deck.NewStore()
	r := api.SetupRouter()
	r.Run(":8080")
}
