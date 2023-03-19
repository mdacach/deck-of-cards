package card

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: These tests could be improved by using property-based testing.
func TestCardString(t *testing.T) {
	testCases := []struct {
		name     string
		card     Card
		expected string
	}{
		{"Ace of Spades -> AS", Card{Suit: Spades, Rank: Ace}, "AS"},
		{"Ten of Hearts -> 10H", Card{Suit: Hearts, Rank: Ten}, "10H"},
		{"Queen of Diamonds -> QD", Card{Suit: Diamonds, Rank: Queen}, "QD"},
		{"Jack of Clubs -> JC", Card{Suit: Clubs, Rank: Jack}, "JC"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := tc.card.String()
			assert.Equal(t, tc.expected, code, "Expected code %q, got: %q", tc.expected, code)
		})
	}
}

func TestCardFromStringValid(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedCard Card
	}{
		{"AS -> Ace of Spades", "AS", Card{Rank: Ace, Suit: Spades}},
		{"KD -> King of Diamonds", "KD", Card{Rank: King, Suit: Diamonds}},
		{"5H -> Five of Hearts", "5H", Card{Rank: Five, Suit: Hearts}},
		{"9C -> Nine of Clubs", "9C", Card{Rank: Nine, Suit: Clubs}},
		{"2S -> Two of Spades", "2S", Card{Rank: Two, Suit: Spades}},
		{"10S -> Ten of Spades", "10S", Card{Rank: Ten, Suit: Spades}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card, err := FromString(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCard, card)
		})
	}
}

// We only check if an error occurred. This way we do not rely on implementation details.
// For example, suppose the input rank and suit are invalid, and in the implementation
// you validate the rank before the suit.
// The error message would be `rank invalid`. (OK)
// If you then change the code to validate the suit before the rank,
// The error message would then be `suit invalid`. (Still OK, but the test would fail)
func TestCardFromStringInvalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"invalid rank string X", "XH"},
		{"invalid suit string $", "A$"},
		{"empty card string", ""},
		{"both rank and suit invalid", "X$"},
		// TODO: UTF-8 handling may require some more investigation. I am not sure how string slices of UTF-8 work in Go.
		{"UTF-8 invalid rank string", "ǶH"},
		{"UTF-8 invalid suit string", "A♡"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FromString(tc.input)
			assert.Error(t, err, "Invalid card code should return an error")
		})
	}
}

func TestCardMarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		card         Card
		expectedJSON string
	}{
		{"Queen of Hearts", Card{Rank: Queen, Suit: Hearts}, `{"value":"QUEEN","suit":"HEARTS","code":"QH"}`},
		{"Ace of Spades", Card{Rank: Ace, Suit: Spades}, `{"value":"ACE","suit":"SPADES","code":"AS"}`},
		{"Ten of Spades", Card{Rank: Ten, Suit: Spades}, `{"value":"TEN","suit":"SPADES","code":"10S"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.card)
			assert.NoError(t, err, "Error marshalling card")
			assert.JSONEq(t, tc.expectedJSON, string(jsonData), "Card JSON representation is not as expected")
		})
	}
}

func TestCardUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		inputJSON    string
		expectedCard Card
	}{
		{"JSON -> Queen of Hearts", `{"value":"QUEEN","suit":"HEARTS","code":"QH"}`, Card{Rank: Queen, Suit: Hearts}},
		{"JSON -> Ace of Spades", `{"value":"ACE","suit":"SPADES","code":"AS"}`, Card{Rank: Ace, Suit: Spades}},
		{"JSON -> Ten of Clubs", `{"value":"TEN","suit":"CLUBS","code":"10"}`, Card{Rank: Ten, Suit: Clubs}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var c Card
			err := json.Unmarshal([]byte(tc.inputJSON), &c)

			assert.NoError(t, err, "Error unmarshalling card")
			assert.Equal(t, tc.expectedCard, c, "Unmarshalled card is not as expected")
		})
	}
}
