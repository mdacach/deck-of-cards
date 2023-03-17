package card

import "fmt"

type Rank string
type Suit string

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

const (
	Spades   Suit = "S"
	Diamonds Suit = "D"
	Clubs    Suit = "C"
	Hearts   Suit = "H"
)

// TODO: Use type-driven-development newtype pattern to make sure string card rank and suit is valid.
type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}
