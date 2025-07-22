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
