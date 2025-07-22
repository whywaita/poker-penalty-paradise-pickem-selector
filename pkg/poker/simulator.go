package poker

// SimulateEquity returns hero win probability against one opponent by Monte‑Carlo
func SimulateEquity(g Game, my4 []Card, iters int) float64 {
	// Special handling for HiDuGi as a split pot game
	if _, ok := g.(HiDuGiGame); ok {
		return SimulateHiDuGiEquity(my4, iters)
	}

	wins, ties := 0, 0
	for i := 0; i < iters; i++ {
		deck := RemoveCards(FullDeck(), ToSet(my4))
		myHand, oppHand, board, _ := g.CompleteHand(my4, deck)
		myScore := g.Evaluate(myHand, board)
		oppScore := g.Evaluate(oppHand, board)
		if myScore > oppScore {
			wins++
		} else if myScore == oppScore {
			ties++
		}
	}
	return (float64(wins) + float64(ties)/2) / float64(iters)
}

// PickBestGame finds the best game variant for the given hand
func PickBestGame(my4 []Card, iters int) (best Game, equities map[string]float64) {
	games := []Game{
		HiDuGiGame{}, // Put HiDuGi first so it wins ties with split pot preference
		DrawmahaHi{},
		BadugiGame{},
		StubGame{"Drawmaha-2-7"},
		StubGame{"Prime"},
		StubGame{"Omaha DoubleBoard"},
	}
	equities = make(map[string]float64, len(games))
	for _, g := range games {
		eq := SimulateEquity(g, my4, iters)
		equities[g.Name()] = eq
		if best == nil || eq > equities[best.Name()] {
			best = g
		}
	}
	return
}
