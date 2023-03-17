package deck

import (
	"deck_of_cards/card"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
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

// TODO: We may want to return a *new* deck here, and not mutate the caller.
// There is no need to have shuffle functionality inside of creating the deck.
// We can first create the deck, then shuffle it (if needed).
func (d *Deck) Shuffle() {
	// TODO: Handle the seed value here. We want to be careful here.
	//       Do not let clients know how the deck is shuffled!
	numberCards := len(d.Cards)
	rand.Shuffle(numberCards, func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw(count int) ([]card.Card, error) {
	if count > d.Remaining {
		return nil, fmt.Errorf("not enough cards remaining in the deck")
	}

	// We copy all drawn cards at once to avoid reallocating the array multiple times.
	drawnCards := make([]card.Card, count)
	copy(drawnCards, d.Cards[:count]) // We don't want to mutate d.Cards from drawnCards.

	// TODO: Write about this.
	//       Here it's not a very big deal to keep resizing the array (because it's small, at most 52 items)
	//       But in another context, we could think about optimizing this: for example, we can keep the array
	//       the same, and keep track of how many cards we have removed until now. Then the "first" card of the
	//       array will actually be at the index `cardsRemoved`.
	d.Cards = d.Cards[count:]
	d.Remaining -= count

	return drawnCards, nil
}
