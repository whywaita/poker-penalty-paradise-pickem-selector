package poker

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
)

// rng is our random number generator
var rng *mrand.Rand

func init() {
	var seed int64
	binary.Read(rand.Reader, binary.LittleEndian, &seed)
	rng = mrand.New(mrand.NewSource(seed))
}

// FullDeck returns a complete 52-card deck
func FullDeck() []Card {
	d := make([]Card, 52)
	for i := 0; i < 52; i++ {
		d[i] = Card(i)
	}
	return d
}

// RemoveCards returns deck with specified cards removed
func RemoveCards(deck []Card, toRemove map[Card]struct{}) []Card {
	out := make([]Card, 0, len(deck))
	for _, c := range deck {
		if _, ok := toRemove[c]; !ok {
			out = append(out, c)
		}
	}
	return out
}

// DrawRandom draws n unique random cards from `deck` in‑place (Fisher‑Yates shuffle prefix)
func DrawRandom(deck []Card, n int) ([]Card, []Card) {
	if n > len(deck) {
		panic("DrawRandom: not enough cards")
	}
	for i := 0; i < n; i++ {
		j := rng.Intn(len(deck)-i) + i
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck[:n], deck[n:]
}

// ToSet converts a slice of cards to a set
func ToSet(cards []Card) map[Card]struct{} {
	m := make(map[Card]struct{}, len(cards))
	for _, c := range cards {
		m[c] = struct{}{}
	}
	return m
}
