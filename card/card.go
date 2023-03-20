// Package card provides types and functions for working with standard playing cards.
// It defines the Card, Rank, and Suit types, along with various utility functions
// for creating and validating cards, as well as converting between short and long
// string representations of ranks and suits.
//
// Example usage:
//
//	c, _ := card.FromString("AS")
//	fmt.Println(c.String()) // Output: AS
//	fmt.Println(c.Rank.LongString()) // Output: ACE
//	fmt.Println(c.Suit.LongString()) // Output: SPADES
package card

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Card represents a single playing card with a rank and suit.
type Card struct {
	Rank Rank
	Suit Suit
}

// FromString creates a Card instance from a string input (e.g., "4H" for the Four of Hearts).
// It returns an error if the input string is invalid.
func FromString(s string) (Card, error) {
	// All codes have at least two characters (one for Rank, one for Suit).
	if len(s) < 2 {
		return Card{}, errors.New("invalid card string")
	}

	// Code is Rank (one char) followed by Suit (one char).
	rankStr := s[0]
	suitStr := s[1]

	// Up until here, we do not know if the rank and suit are valid strings.
	// So we need to validate them.
	rank := Rank(rankStr)
	suit := Suit(suitStr)

	if !rank.IsValid() {
		return Card{}, fmt.Errorf("invalid rank string: %c", rankStr)
	}

	if !suit.IsValid() {
		return Card{}, fmt.Errorf("invalid suit string: %c", suitStr)
	}

	return Card{Rank: rank, Suit: suit}, nil
}

// MarshalJSON customizes the JSON marshaling of the Card struct. It returns a JSON object
// with the long form value of the rank, long form value of the suit, and the card code.
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

// UnmarshalJSON customizes the JSON unmarshalling of the Card struct. It expects a JSON object
// with the long form value of the rank and the long form value of the suit. It returns an error
// if the input JSON is invalid.
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

// String returns a string representation of the Card (e.g., "4H" for the Four of Hearts).
// This is also referred as the `Code` of the Card.
func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}
