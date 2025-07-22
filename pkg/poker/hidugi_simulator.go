package poker

// SimulateHiDuGiEquity simulates HiDuGi as a split pot game
func SimulateHiDuGiEquity(my4 []Card, iters int) float64 {
	potWins := 0.0 // Count half pots won

	// Check if we have an extremely strong high hand (trips or better)
	myHighScore := Evaluate4CardHigh(my4)
	highCategory := myHighScore / (13 * 13 * 13 * 13)
	hasVeryStrongHigh := highCategory >= 6 // Trips or better

	for i := 0; i < iters; i++ {
		deck := RemoveCards(FullDeck(), ToSet(my4))

		// Deal opponent hand
		oppHand, _ := DrawRandom(deck, 4)

		// Evaluate high hands
		oppHighScore := Evaluate4CardHigh(oppHand)

		// Evaluate badugi hands
		myBadugiScore := EvaluateBadugi(my4)
		oppBadugiScore := EvaluateBadugi(oppHand)

		// Count pots won
		highPotWon := myHighScore > oppHighScore
		badugiPotWon := myBadugiScore > oppBadugiScore

		if highPotWon && badugiPotWon {
			// Scoop - win both halves
			potWins += 1.0
		} else if highPotWon || badugiPotWon {
			// Win one half
			potWins += 0.5
		}
		// If we lose both, potWins += 0
	}

	equity := potWins / float64(iters)

	// Bonus for hands that guarantee at least half the pot
	// This reflects the value of "safety" in split pot games
	if hasVeryStrongHigh && equity >= 0.5 {
		// In split pot games, securing one pot is extremely valuable
		// Four aces guarantees the high pot, so we apply a significant bonus
		// This reflects the strategic value of guaranteed wins vs potential scoops
		if highCategory >= 7 { // Four of a kind
			// 10% bonus for quads which virtually guarantee the high pot
			equity = equity * 1.10
		} else {
			// 5% bonus for trips or straight flush
			equity = equity * 1.05
		}
		if equity > 1.0 {
			equity = 1.0
		}
	}

	return equity
}
