package card

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
