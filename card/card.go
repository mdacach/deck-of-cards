package card

import (
	"encoding/json"
	"fmt"
)

// TODO: Use type-driven-development newtype pattern to make sure string card rank and suit is valid.
type Card struct {
	Rank Rank
	Suit Suit
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

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}
