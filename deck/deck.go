package deck

import (
	"deck_of_cards/card"
	"fmt"
	"github.com/google/uuid"
)

type Deck struct {
	ID        uuid.UUID
	Cards     []card.Card
	Remaining int
}

func NewStandardDeck() Deck {
	cards := make([]card.Card, 0, 52)

	for _, s := range card.Suits() {
		for _, r := range card.Ranks() {
			cards = append(cards, card.Card{Rank: r, Suit: s})
		}
	}

	return Deck{
		ID:        uuid.New(),
		Cards:     cards,
		Remaining: len(cards),
	}
}

// deck.go

func NewPartialDeck(codes []string) (Deck, error) {
	cards := make([]card.Card, 0, len(codes))

	for _, code := range codes {
		c, err := card.FromString(code)
		if err != nil {
			return Deck{}, fmt.Errorf("invalid card code '%s': %w", code, err)
		}
		cards = append(cards, c)
	}

	return Deck{
		ID:        uuid.New(),
		Cards:     cards,
		Remaining: len(cards),
	}, nil
}
