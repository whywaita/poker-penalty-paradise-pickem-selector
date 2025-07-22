package poker

// Game interface defines poker game variants
type Game interface {
	Name() string
	// CompleteHand fills in missing private and public cards for simulation.
	CompleteHand(my []Card, deck []Card) (myComplete []Card, oppHand []Card, board []Card, deckOut []Card)
	Evaluate(myComplete []Card, board []Card) int64
}

// DrawmahaHi implementation
type DrawmahaHi struct{}

func (d DrawmahaHi) Name() string { return "Drawmaha-Hi" }

func (d DrawmahaHi) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	// Drawmaha deals 5‑card hands, no community board. Each player keeps 4 original cards & draws 1.
	var myHand, oppHand []Card
	var drawn []Card
	// Draw 1 card for hero
	drawn, deck = DrawRandom(deck, 1)
	myHand = append(append([]Card(nil), my...), drawn...)
	// Opponent 5 cards
	oppHand, deck = DrawRandom(deck, 5)
	return myHand, oppHand, nil, deck
}

func (d DrawmahaHi) Evaluate(h []Card, board []Card) int64 {
	// Already completed to 5 cards
	return Evaluate5CardHigh(h)
}

// BadugiGame implementation
type BadugiGame struct{}

func (b BadugiGame) Name() string { return "Badugi" }

func (b BadugiGame) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	// Badugi uses 4‑card hands; hero already has 4.
	oppHand, deck := DrawRandom(deck, 4)
	return my, oppHand, nil, deck
}

func (b BadugiGame) Evaluate(h []Card, board []Card) int64 { return EvaluateBadugi(h) }

// HiDuGiGame implementation - split pot Hi/Badugi game
type HiDuGiGame struct{}

func (h HiDuGiGame) Name() string { return "HiDuGi" }

func (h HiDuGiGame) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	// HiDuGi uses 4-card hands; hero already has 4.
	oppHand, deck := DrawRandom(deck, 4)
	return my, oppHand, nil, deck
}

func (h HiDuGiGame) Evaluate(hand []Card, board []Card) int64 {
	highScore, badugiScore := EvaluateHiDuGi(hand)

	// Check if we have 8-badugi or better for a strong hand
	has8Badugi := IsBadugi8OrBetter(hand)

	// Scoring strategy:
	// - If we have 8-badugi or better, heavily weight the combined score
	// - Otherwise, use a balanced approach between high and badugi

	// Get the high hand category (0-8)
	highCategory := highScore / (13 * 13 * 13 * 13)

	// Normalize badugi score
	normalizedBadugi := badugiScore / 1000000000000

	if has8Badugi {
		// Strong badugi hands get significant bonus
		// Add bonus of 10 million to ensure these hands rank highest
		return highScore + normalizedBadugi + 10000000
	}

	// For hands without 8-badugi:
	// - Strong high hands (category 6+: trips, quads, straight flush) get a boost
	// - This ensures four aces (category 7) beats most other hands
	if highCategory >= 6 {
		// Trips or better get a significant boost
		return highScore*3 + normalizedBadugi/10
	}

	// For other hands, combine scores with equal weight
	return highScore + normalizedBadugi/10
}

// StubGame implementation for unimplemented variants
type StubGame struct {
	NameStr string
}

func (s StubGame) Name() string { return s.NameStr }
func (s StubGame) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	oppHand, deck := DrawRandom(deck, len(my))
	return my, oppHand, nil, deck
}
func (s StubGame) Evaluate(h []Card, board []Card) int64 { return 0 }
