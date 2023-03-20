package card

import (
	"fmt"
	"strings"
)

// Rank represents the rank of a playing card (e.g., Ace, Two, ... King).
type Rank string

const (
	Ace   Rank = "A"
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "T" // Poker uses "T" instead of Ten. [TJ (Ten-Jack) suited]
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
)

// Ranks returns a slice of all valid Rank values, in order (Ace first).
func Ranks() []Rank {
	return []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
}

// IsValid checks whether the Rank is a valid value.
func (r Rank) IsValid() bool {
	for _, validRank := range Ranks() {
		if r == validRank {
			return true
		}
	}
	return false
}

// LongString returns the long form string representation of the Rank (e.g., "ACE", "TWO", ... "KING").
func (r Rank) LongString() string {
	switch r {
	case Ace:
		return "ACE"
	case Two:
		return "TWO"
	case Three:
		return "THREE"
	case Four:
		return "FOUR"
	case Five:
		return "FIVE"
	case Six:
		return "SIX"
	case Seven:
		return "SEVEN"
	case Eight:
		return "EIGHT"
	case Nine:
		return "NINE"
	case Ten:
		return "TEN"
	case Jack:
		return "JACK"
	case Queen:
		return "QUEEN"
	case King:
		return "KING"
	default:
		return ""
	}
}

// ParseLongRank takes a long form rank string and returns the corresponding Rank value.
// It returns an error if the input string is not a valid long form rank.
func ParseLongRank(r string) (Rank, error) {
	switch strings.ToUpper(r) {
	case "ACE":
		return Ace, nil
	case "TWO":
		return Two, nil
	case "THREE":
		return Three, nil
	case "FOUR":
		return Four, nil
	case "FIVE":
		return Five, nil
	case "SIX":
		return Six, nil
	case "SEVEN":
		return Seven, nil
	case "EIGHT":
		return Eight, nil
	case "NINE":
		return Nine, nil
	case "TEN":
		return Ten, nil
	case "JACK":
		return Jack, nil
	case "QUEEN":
		return Queen, nil
	case "KING":
		return King, nil
	default:
		return "", fmt.Errorf("could not parse Rank from string: %s", r)
	}
}
