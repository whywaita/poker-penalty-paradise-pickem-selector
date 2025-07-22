package poker

import (
	"math"
)

// Hand category constants (higher is better)
const (
	HighCard = iota
	OnePair
	TwoPair
	Trips
	Straight
	Flush
	FullHouse
	Quads
	StraightFlush
)

// Evaluate5CardHigh evaluates a 5-card poker hand - category ranking (higher is better)
func Evaluate5CardHigh(hand []Card) int64 {
	if len(hand) != 5 {
		panic("evaluate5CardHigh expects 5 cards")
	}
	ranks := make([]int, 13)
	suits := make([]int, 4)
	for _, c := range hand {
		ranks[c.Rank()]++
		suits[c.Suit()]++
	}
	isFlush := false
	for _, s := range suits {
		if s == 5 {
			isFlush = true
			break
		}
	}
	// detect straight (wheel A2345)
	top := -1
	consec := 0
	for r := 12; r >= 0; r-- {
		if ranks[r] > 0 {
			consec++
			if consec == 5 {
				top = r + 4 // highest rank of straight
				break
			}
		} else {
			consec = 0
		}
	}
	// wheel check
	if consec == 4 && ranks[12] > 0 && ranks[3] > 0 && ranks[2] > 0 && ranks[1] > 0 && ranks[0] > 0 {
		top = 3 // straight to 5
	}
	isStraight := top != -1

	// count frequencies
	var pairs, trips, quads int
	for _, cnt := range ranks {
		switch cnt {
		case 2:
			pairs++
		case 3:
			trips++
		case 4:
			quads++
		}
	}

	var cat int64
	switch {
	case isStraight && isFlush:
		cat = StraightFlush
	case quads > 0:
		cat = Quads
	case trips > 0 && pairs > 0:
		cat = FullHouse
	case isFlush:
		cat = Flush
	case isStraight:
		cat = Straight
	case trips > 0:
		cat = Trips
	case pairs >= 2:
		cat = TwoPair
	case pairs > 0:
		cat = OnePair
	default:
		cat = HighCard
	}

	// naive kicker encoding: combine sorted ranks descending into value
	vals := make([]int, 0, 5)
	for r := 12; r >= 0; r-- {
		for i := 0; i < ranks[r]; i++ {
			vals = append(vals, r)
		}
	}
	var kicker int64
	for _, v := range vals {
		kicker = kicker*13 + int64(v)
	}
	return cat*int64(13*13*13*13*13) + kicker // category has highest weight
}

// EvaluateBadugi evaluates a Badugi hand – lower is better (4‑card low with unique suits)
// We convert to HIGH by negating within a large offset so that higher = better overall
func EvaluateBadugi(hand []Card) int64 {
	// select up to 4 cards with distinct suits and ranks minimal lexicographically.
	chosen := make([]Card, 0, 4)
	usedSuit := [4]bool{}
	usedRank := [13]bool{}
	// brute force: iterate ranks ascending, pick first available card with unique suit.
	for r := 0; r < 13 && len(chosen) < 4; r++ {
		for s := 0; s < 4 && len(chosen) < 4; s++ {
			c := Card(r*4 + s)
			for _, h := range hand {
				if h == c && !usedSuit[s] && !usedRank[r] {
					chosen = append(chosen, h)
					usedSuit[s] = true
					usedRank[r] = true
					break
				}
			}
		}
	}
	// score: more cards better; then lower ranks better.
	// Base score by number of cards (4 cards = 4000, 3 cards = 3000, etc)
	score := int64(4-len(chosen)) * 100000
	// Add rank values (lower is better in badugi)
	for _, c := range chosen {
		score = score*13 + int64(c.Rank())
	}
	// invert so bigger is better (align with evaluate5CardHigh)
	return math.MaxInt64/2 - score
}

// Evaluate4CardHigh evaluates a 4-card poker hand - category ranking (higher is better)
func Evaluate4CardHigh(hand []Card) int64 {
	if len(hand) != 4 {
		panic("evaluate4CardHigh expects 4 cards")
	}
	ranks := make([]int, 13)
	suits := make([]int, 4)
	for _, c := range hand {
		ranks[c.Rank()]++
		suits[c.Suit()]++
	}

	// Check for flush (4 cards of same suit)
	isFlush := false
	for _, s := range suits {
		if s == 4 {
			isFlush = true
			break
		}
	}

	// Check for straight (4 consecutive cards)
	isStraight := false
	consec := 0
	for r := 12; r >= 0; r-- {
		if ranks[r] > 0 {
			consec++
			if consec == 4 {
				isStraight = true
				break
			}
		} else {
			consec = 0
		}
	}
	// Check wheel (A234)
	if !isStraight && ranks[12] > 0 && ranks[0] > 0 && ranks[1] > 0 && ranks[2] > 0 {
		isStraight = true
	}

	// count frequencies
	var pairs, trips, quads int
	for _, cnt := range ranks {
		switch cnt {
		case 2:
			pairs++
		case 3:
			trips++
		case 4:
			quads++
		}
	}

	var cat int64
	switch {
	case isStraight && isFlush:
		cat = 8 // Straight flush
	case quads > 0:
		cat = 7 // Four of a kind
	case trips > 0:
		cat = 6 // Three of a kind
	case isFlush:
		cat = 5 // Flush
	case isStraight:
		cat = 4 // Straight
	case pairs >= 2:
		cat = 3 // Two pair
	case pairs > 0:
		cat = 2 // One pair
	default:
		cat = 1 // High card
	}

	// Kicker encoding for 4 cards
	vals := make([]int, 0, 4)
	for r := 12; r >= 0; r-- {
		for i := 0; i < ranks[r]; i++ {
			vals = append(vals, r)
		}
	}
	var kicker int64
	for _, v := range vals {
		kicker = kicker*13 + int64(v)
	}
	return cat*int64(13*13*13*13) + kicker
}

// EvaluateHiDuGi evaluates both high and badugi hands for HiDuGi split pot game
func EvaluateHiDuGi(hand []Card) (int64, int64) {
	highScore := Evaluate4CardHigh(hand)
	badugiScore := EvaluateBadugi(hand)
	return highScore, badugiScore
}

// IsBadugi8OrBetter checks if the hand qualifies as 8-badugi or better
func IsBadugi8OrBetter(hand []Card) bool {
	// Get the badugi hand
	chosen := make([]Card, 0, 4)
	usedSuit := [4]bool{}
	usedRank := [13]bool{}

	for r := 0; r < 13 && len(chosen) < 4; r++ {
		for s := 0; s < 4 && len(chosen) < 4; s++ {
			c := Card(r*4 + s)
			for _, h := range hand {
				if h == c && !usedSuit[s] && !usedRank[r] {
					chosen = append(chosen, h)
					usedSuit[s] = true
					usedRank[r] = true
					break
				}
			}
		}
	}

	// Check if it's a 4-card badugi with highest card 8 or better
	// In badugi, ace is low, so we need to check the actual high card
	if len(chosen) == 4 {
		maxRank := -1
		for _, c := range chosen {
			rank := c.Rank()
			// Ace is low in badugi, treat as rank 0 for comparison
			if rank == 12 {
				rank = -1
			}
			if rank > maxRank {
				maxRank = rank
			}
		}
		// maxRank 6 = 8, so <= 6 means 8 or better
		return maxRank <= 6
	}
	return false
}
