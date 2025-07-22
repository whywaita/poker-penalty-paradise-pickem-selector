package poker

import (
	"fmt"
	"strings"
)

// Card is an int 0‑51 where rank = c/4 (0=2 .. 12=Ace) and suit = c%4 (0=♣,1=♦,2=♥,3=♠)
// This mapping makes bit‑operations convenient.
type Card int

var rankToChar = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
var suitToChar = []string{"c", "d", "h", "s"}

// CardFromString creates a card from string representation (e.g., "As", "2c")
func CardFromString(s string) (Card, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid card string %q", s)
	}
	r, su := strings.ToUpper(string(s[0])), strings.ToLower(string(s[1]))
	var rank, suit int = -1, -1
	for i, c := range rankToChar {
		if strings.EqualFold(c, r) {
			rank = i
			break
		}
	}
	for i, c := range suitToChar {
		if c == su {
			suit = i
			break
		}
	}
	if rank == -1 || suit == -1 {
		return 0, fmt.Errorf("invalid card %q", s)
	}
	return Card(rank*4 + suit), nil
}

// Rank returns the rank of the card (0=2 .. 12=Ace)
func (c Card) Rank() int { return int(c) / 4 }

// Suit returns the suit of the card (0=♣,1=♦,2=♥,3=♠)
func (c Card) Suit() int { return int(c) % 4 }

// String returns the string representation of the card
func (c Card) String() string {
	return rankToChar[c.Rank()] + suitToChar[c.Suit()]
}
