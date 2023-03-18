package deck

import (
	"deck_of_cards/card"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

// Deck represents a deck of playing cards.
type Deck struct {
	// ID is a unique identifier for the deck.
	ID uuid.UUID
	// Shuffled indicates whether the deck has been shuffled or not.
	Shuffled bool
	// Remaining represents the number of cards remaining to be drawn in the deck.
	Remaining int
	// Cards holds the card objects in the deck.
	// Cards are specified in draw-order (the first one in the array will be drawn first).
	Cards []card.Card
}

// NewStandardDeck creates a new Deck containing a full set of 52 standard playing cards.
func NewStandardDeck() Deck {
	cards := make([]card.Card, 0, 52)

	for _, s := range card.Suits() {
		for _, r := range card.Ranks() {
			cards = append(cards, card.Card{Rank: r, Suit: s})
		}
	}

	return Deck{
		ID:        uuid.New(),
		Shuffled:  false,
		Remaining: len(cards),
		Cards:     cards,
	}
}

// NewPartialDeck creates a new Deck containing a custom set of cards based on the provided card codes.
// It returns an error if any of these happens:
// 1. The codes array is empty.
// 2. There are any invalid codes in the codes array.
// 3. There are repeated codes in the codes array.
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
		Shuffled:  false,
		Remaining: len(cards),
		Cards:     cards,
	}, nil
}

// Shuffle shuffles the cards in the Deck. Note that this mutates the Deck.
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

	d.Shuffled = true
}

// Draw removes and returns the specified number of cards from the top (the front) of the Deck.
// It returns an error if there are not enough cards remaining in the Deck.
func (d *Deck) Draw(count int) ([]card.Card, error) {
	if count > d.Remaining {
		return nil, fmt.Errorf("not enough cards remaining in the deck")
	}

	if count <= 0 {
		return nil, fmt.Errorf("draw count should be positive")
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
