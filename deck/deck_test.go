package deck

import (
	"deck_of_cards/card"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStandardDeck(t *testing.T) {
	deck := NewStandardDeck()

	assert.NotNil(t, deck.ID, "Deck ID should not be nil")
	assert.Len(t, deck.Cards, 52, "A standard deck should have 52 cards")
	assert.Equal(t, deck.Remaining, len(deck.Cards), "A standard deck should start with 52 cards remaining")

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

func TestNewPartialDeckCards(t *testing.T) {
	codes := []string{"AS", "KD", "AC", "2C", "KH"}

	deck, err := NewPartialDeck(codes)

	assert.NoError(t, err, "Creating a partial deck with example codes should not return an error")
	// Relevant fields are populated.
	assert.NotNil(t, deck.ID, "Deck ID should not be nil")
	assert.Len(t, deck.Cards, len(codes), "Partial deck should have the specified number of cards")
	assert.Equal(t, len(codes), deck.Remaining, "Remaining cards should match the number of input codes")

	// The cards we passed by code are the deck's cards.
	for i, code := range codes {
		assert.Equal(t, code, deck.Cards[i].String(), "Deck card should match the specified code")
	}
}

func TestNewPartialDeckValid(t *testing.T) {
	testCases := []struct {
		name        string
		cardStrings []string
		wantDeckLen int
	}{
		{
			name:        "single valid card",
			cardStrings: []string{"AS"},
			wantDeckLen: 1,
		},
		{
			name:        "multiple valid cards - no repeated",
			cardStrings: []string{"AS", "KD", "AC", "2C", "KH"},
			wantDeckLen: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deck, err := NewPartialDeck(tc.cardStrings)
			// If there is an error, we do not want to continue with the execution.
			require.NoError(t, err, "Unexpected error in NewPartialDeck: %v", err)
			assert.Equal(t, tc.wantDeckLen, len(deck.Cards), "Expected deck length to be: %v but got: %v", tc.wantDeckLen, len(deck.Cards))
			assert.Equal(t, tc.wantDeckLen, deck.Remaining, "Expected deck remaining cards to be: %v but got: %v", tc.wantDeckLen, deck.Remaining)
		})
	}
}
