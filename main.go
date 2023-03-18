package main

import "deck_of_cards/api"

func main() {
	r := api.SetupRouter()
	r.Run(":8080")
}
