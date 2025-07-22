package poker

import (
	"fmt"
	"strings"
)

// ParseHand parses a hand string into cards
func ParseHand(arg string) ([]Card, error) {
	parts := strings.Fields(arg)
	if len(parts) == 1 {
		// maybe comma‑separated "Ac,Kd,2h,3c"
		parts = strings.Split(arg, ",")
	}
	if len(parts) != 4 {
		return nil, fmt.Errorf("need exactly 4 cards, got %d", len(parts))
	}
	hand := make([]Card, 0, 4)
	seen := map[Card]struct{}{}
	for _, p := range parts {
		c, err := CardFromString(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		if _, dup := seen[c]; dup {
			return nil, fmt.Errorf("duplicate card %s", c)
		}
		seen[c] = struct{}{}
		hand = append(hand, c)
	}
	return hand, nil
}
