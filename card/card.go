package card

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TODO: Use type-driven-development newtype pattern to make sure string card rank and suit is valid.
type Card struct {
	Rank Rank
	Suit Suit
}

func FromString(s string) (Card, error) {
	// All codes have at least two characters (one for Rank, one for Suit).
	if len(s) < 2 {
		return Card{}, errors.New("invalid card string")
	}

	// The suit should come always last, and have always length one.
	// So we can split the string into the rank and suit portions accordingly.
	rankStr := s[:len(s)-1]
	suitStr := s[len(s)-1:]

	// Up until here, we do not know if the rank and suit are valid strings.
	// So we need to validate them.
	// TODO: This will change.
	rank := Rank(rankStr)
	suit := Suit(suitStr)

	if !rank.IsValid() {
		return Card{}, fmt.Errorf("invalid rank string: %s", rankStr)
	}

	if !suit.IsValid() {
		return Card{}, fmt.Errorf("invalid suit string: %s", suitStr)
	}

	return Card{Rank: rank, Suit: suit}, nil
}

func (c Card) MarshalJSON() ([]byte, error) {
	cardJSON := struct {
		Value string `json:"value"`
		Suit  string `json:"suit"`
		Code  string `json:"code"`
	}{
		Value: c.Rank.LongString(),
		Suit:  c.Suit.LongString(),
		Code:  c.String(),
	}

	return json.Marshal(cardJSON)
}

func (c *Card) UnmarshalJSON(data []byte) error {
	cardJSON := struct {
		Value string `json:"value"`
		Suit  string `json:"suit"`
	}{}

	if err := json.Unmarshal(data, &cardJSON); err != nil {
		return err
	}

	rank, err := ParseLongRank(cardJSON.Value)
	if err != nil {
		return err
	}

	suit, err := ParseLongSuit(cardJSON.Suit)
	if err != nil {
		return err
	}

	c.Rank = rank
	c.Suit = suit

	return nil
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}
