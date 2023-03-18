package card

import (
	"fmt"
	"strings"
)

// Suit represents the suit of a playing card (Spades, Diamonds, Clubs, Hearts).
type Suit string

const (
	Spades   Suit = "S"
	Diamonds Suit = "D"
	Clubs    Suit = "C"
	Hearts   Suit = "H"
)

// Suits returns a slice of all valid Suit values, in order (Spades, Diamonds, Clubs, Hearts).
func Suits() []Suit {
	return []Suit{Spades, Diamonds, Clubs, Hearts}
}

// IsValid checks whether the Suit is a valid value.
// TODO: This will be refactored. The validation will happen in the Rank constructor.
func (s Suit) IsValid() bool {
	for _, validSuit := range Suits() {
		if s == validSuit {
			return true
		}
	}
	return false
}

// LongString returns the long form string representation of the Suit (e.g., "SPADES", "DIAMONDS", "CLUBS", "HEARTS").
func (s Suit) LongString() string {
	switch s {
	case Spades:
		return "SPADES"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	case Hearts:
		return "HEARTS"
	default:
		// TODO: Validate that Suit is always a valid variant (by construction).
		return ""
	}
}

// ParseLongSuit takes a long form suit string and returns the corresponding Suit value.
// It returns an error if the input string is not a valid long form suit.
func ParseLongSuit(s string) (Suit, error) {
	switch strings.ToUpper(s) {
	case "SPADES":
		return Spades, nil
	case "DIAMONDS":
		return Diamonds, nil
	case "CLUBS":
		return Clubs, nil
	case "HEARTS":
		return Hearts, nil
	default:
		return "", fmt.Errorf("could not parse Suit from string: %s", s)
	}
}
