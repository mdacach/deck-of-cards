package deck

import (
	"deck_of_cards/card"
	"errors"
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

func NewPartialDeck(codes []string) (Deck, error) {
	if len(codes) == 0 {
		return Deck{}, errors.New("a deck must have at least one card")
	}

	cards := make([]card.Card, 0, len(codes))

	for _, code := range codes {
		c, err := card.FromString(code)
		if err != nil {
			return Deck{}, fmt.Errorf("invalid card code '%s': %w", code, err)
		}
		cards = append(cards, c)
	}

	cardSet := make(map[string]bool)
	for _, code := range codes {
		if _, exists := cardSet[code]; exists {
			return Deck{}, errors.New("repeated card code")
		}
		cardSet[code] = true
	}

	return Deck{
		ID:        uuid.New(),
		Cards:     cards,
		Remaining: len(cards),
	}, nil
}
