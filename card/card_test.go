package card

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: These tests could be improved by using property-based testing.
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

func TestCardFromStringValid(t *testing.T) {
	testCases := []struct {
		input        string
		expectedCard Card
	}{
		{"AS", Card{Rank: Ace, Suit: Spades}},
		{"KD", Card{Rank: King, Suit: Diamonds}},
		{"5H", Card{Rank: Five, Suit: Hearts}},
		{"9C", Card{Rank: Nine, Suit: Clubs}},
		{"2S", Card{Rank: Two, Suit: Spades}},
		{"10S", Card{Rank: Ten, Suit: Spades}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("FromString_%s", tc.input), func(t *testing.T) {
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
			t.Logf("error: %s", err)
			assert.Error(t, err)
		})
	}
}

func TestCardMarshalJSON(t *testing.T) {
	testCases := []struct {
		card         Card
		expectedJSON string
	}{
		{Card{Rank: Queen, Suit: Hearts}, `{"value":"QUEEN","suit":"HEARTS","code":"QH"}`},
		{Card{Rank: Ace, Suit: Spades}, `{"value":"ACE","suit":"SPADES","code":"AS"}`},
	}

	for _, tc := range testCases {
		jsonData, err := json.Marshal(tc.card)
		assert.NoError(t, err, "Error marshaling card")
		assert.JSONEq(t, tc.expectedJSON, string(jsonData), "Card JSON representation is not as expected")
	}
}

func TestCardUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		inputJSON    string
		expectedCard Card
	}{
		{`{"value":"QUEEN","suit":"HEARTS","code":"QH"}`, Card{Rank: Queen, Suit: Hearts}},
		{`{"value":"ACE","suit":"SPADES","code":"AS"}`, Card{Rank: Ace, Suit: Spades}},
	}

	for _, tc := range testCases {
		var c Card
		err := json.Unmarshal([]byte(tc.inputJSON), &c)
		assert.NoError(t, err, "Error unmarshalling card")
		assert.Equal(t, tc.expectedCard, c, "Unmarshalled card is not as expected")
	}
}
