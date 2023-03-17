package card_test

import (
	"deck_of_cards/card"
	"testing"
)

func TestCardString(t *testing.T) {
	cases := []struct {
		card     card.Card
		expected string
	}{
		{card: card.Card{Suit: card.Spades, Rank: card.Ace}, expected: "AS"},
		{card: card.Card{Suit: card.Hearts, Rank: card.Ten}, expected: "10H"},
		{card: card.Card{Suit: card.Diamonds, Rank: card.Queen}, expected: "QD"},
		{card: card.Card{Suit: card.Clubs, Rank: card.Jack}, expected: "JC"},
	}

	for _, c := range cases {
		code := c.card.String()
		if code != c.expected {
			t.Errorf("Expected card code %q, but got %q", c.expected, code)
		}
	}
}
