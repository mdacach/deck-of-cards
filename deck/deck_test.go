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
	assert.Equal(t, deck.Shuffled, false, "A standard deck is not shuffled by default")
	assert.Equal(t, deck.Remaining, len(deck.Cards), "A standard deck should start with 52 cards remaining")
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

func TestNewPartialDeckCards(t *testing.T) {
	codes := []string{"AS", "KD", "AC", "2C", "KH"}

	deck, err := NewPartialDeck(codes)

	assert.NoError(t, err, "Creating a partial deck with example codes should not return an error")
	// Relevant fields are populated.
	assert.NotNil(t, deck.ID, "Deck ID should not be nil")
	assert.Equal(t, deck.Shuffled, false, "A partial deck is not shuffled by default")
	assert.Equal(t, len(codes), deck.Remaining, "Remaining cards should match the number of input codes")
	assert.Len(t, deck.Cards, len(codes), "Partial deck should have the specified number of cards")

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

func TestNewPartialDeckInvalidScenarios(t *testing.T) {
	testCases := []struct {
		name        string
		cardStrings []string
	}{
		{
			name:        "empty deck",
			cardStrings: []string{},
		},
		{
			name:        "single invalid codes",
			cardStrings: []string{"ZD"},
		},
		{
			name:        "invalid card code",
			cardStrings: []string{"AS", "ZD", "AC", "2C", "KH"},
		},
		{
			name:        "multiple invalid card codes",
			cardStrings: []string{"AS", "ZD", "AC", "ZZ", "2C", "KH"},
		},
		{
			name:        "only invalid codes",
			cardStrings: []string{"AA", "BB", "33"},
		},
		{
			name:        "invalid utf8 code",
			cardStrings: []string{"AS", "😀D", "AC", "2C", "KH"},
		},
		{
			name:        "repeated card codes",
			cardStrings: []string{"AS", "AS", "AC", "2C", "KH"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPartialDeck(tc.cardStrings)
			require.Error(t, err, "expected an error for %s", tc.name)
		})
	}
}

func TestShuffle(t *testing.T) {
	createPartialDeck := func() Deck { d, _ := NewPartialDeck([]string{"AS", "KD", "AC", "2C", "KH"}); return d }
	testCases := []struct {
		name       string
		createDeck func() Deck
	}{
		{
			name:       "full deck",
			createDeck: NewStandardDeck,
		},
		{
			name:       "partial deck",
			createDeck: createPartialDeck,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := tc.createDeck()

			originalCards := make([]card.Card, len(d.Cards))
			copy(originalCards, d.Cards)

			d.Shuffle()

			assert.Equal(t, len(originalCards), len(d.Cards), "the number of cards should remain the same after shuffling")
			// There is a *very* small probability of this test failing (the shuffle may end up with the cards in the same place).
			// Sorry if that happens to you...
			assert.NotEqual(t, originalCards, d.Cards, "the order of cards should change after shuffling")
			assert.Equal(t, d.Shuffled, true, "deck is marked shuffled")
		})
	}
}

func TestDeckDraw(t *testing.T) {
	testCases := []struct {
		name          string
		drawCount     int
		expectedErr   bool
		expectedCards int
	}{
		{
			name:        "draw zero cards",
			drawCount:   0,
			expectedErr: true,
		},
		{
			name:          "draw one card",
			drawCount:     1,
			expectedErr:   false,
			expectedCards: 1,
		},
		{
			name:          "draw five cards",
			drawCount:     5,
			expectedErr:   false,
			expectedCards: 5,
		},
		{
			name:          "draw exactly the remaining cards",
			drawCount:     52,
			expectedErr:   false,
			expectedCards: 52,
		},
		{
			name:        "draw more cards than remaining",
			drawCount:   60,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deck := NewStandardDeck()

			cards, err := deck.Draw(tc.drawCount)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, cards, tc.expectedCards)
			}
		})
	}
}

func TestDeckDrawOrder(t *testing.T) {
	createPartialDeck := func() Deck { d, _ := NewPartialDeck([]string{"AS", "KD", "AC", "2C", "KH"}); return d }
	testCases := []struct {
		name       string
		shuffled   bool
		drawCount  int
		createDeck func() Deck
	}{
		{
			name:       "draw 3 cards",
			shuffled:   true,
			drawCount:  3,
			createDeck: createPartialDeck,
		},
		{
			name:       "draw 5 cards",
			shuffled:   false,
			drawCount:  5,
			createDeck: NewStandardDeck,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deck := tc.createDeck()
			if tc.shuffled {
				deck.Shuffle()
			}

			orderedCards := deck.Cards
			firstCountCards := orderedCards[:tc.drawCount]

			drawnCards, err := deck.Draw(tc.drawCount)

			require.NoError(t, err)
			require.Len(t, drawnCards, tc.drawCount)

			assert.Equal(t, drawnCards, firstCountCards)
		})
	}
}
