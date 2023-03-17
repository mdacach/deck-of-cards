package card

import (
	"fmt"
	"strings"
)

type Suit string

const (
	Spades   Suit = "S"
	Diamonds Suit = "D"
	Clubs    Suit = "C"
	Hearts   Suit = "H"
)

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
