package card

import "fmt"

// TODO: Use type-driven-development newtype pattern to make sure string card rank and suit is valid.
type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}
