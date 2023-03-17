package deck

import (
	"deck_of_cards/card"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStandardDeck(t *testing.T) {
	deck := NewStandardDeck()

	assert.NotNil(t, deck.ID, "Deck ID should not be nil")
	assert.Len(t, deck.Cards, 52, "A standard deck should have 52 cards")

	cardCount := make(map[card.Card]int)

	for _, c := range deck.Cards {
		cardCount[c]++
	}

	for _, s := range card.Suits() {
		for _, r := range card.Ranks() {
			c := card.Card{Rank: r, Suit: s}
			assert.Equal(t, 1, cardCount[c], "There should be exactly one of each card in the deck")
		}
	}
}
