package poker

func SimulateBidugiEquity(my4 []Card, iters int) float64 {
	potWins := 0.0

	for i := 0; i < iters; i++ {
		deck := RemoveCards(FullDeck(), ToSet(my4))

		// Deal opponent hand
		opp4, _ := DrawRandom(deck, 4)

		my4Final, opp4Final := changeForBadugi(my4, opp4)

		rankMy4Final := getBadugiRank(my4Final)
		rankOpp4Final := getBadugiRank(opp4Final)

		if rankMy4Final < rankOpp4Final {
			potWins += 1.0
		} else if rankMy4Final == rankOpp4Final {
			potWins += 0.5
		}
	}

	equity := potWins / float64(iters)
	return equity
}
