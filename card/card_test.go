package card

import (
	"testing"
)

func TestCardString(t *testing.T) {
	cases := []struct {
		card     Card
		expected string
	}{
		{card: Card{Suit: Spades, Rank: Ace}, expected: "AS"},
		{card: Card{Suit: Hearts, Rank: Ten}, expected: "10H"},
		{card: Card{Suit: Diamonds, Rank: Queen}, expected: "QD"},
		{card: Card{Suit: Clubs, Rank: Jack}, expected: "JC"},
	}

	for _, test := range cases {
		code := test.card.String()
		if code != test.expected {
			t.Errorf("Expected   code %q, but got %q", test.expected, code)
		}
	}
}
