package card

import (
	"fmt"
	"strings"
)

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
	Ten   Rank = "10"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
)

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
		// TODO: Should we return an error instead?
		//       In Rust I would return a Result here, but maybe considering the empty string
		//       as an error in Go is simpler and accomplishes the same as an err variable.
		return ""
	}
}

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
