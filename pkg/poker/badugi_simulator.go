package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func SimulateBadugiEquity(my4 []Card, iters int) float64 {
	potWins := 0.0

	for i := 0; i < iters; i++ {
		deck := RemoveCards(FullDeck(), ToSet(my4))

		// Deal opponent hand
		opp4, _ := DrawRandom(deck, 4)

		my4Final, opp4Final := changeForBadugi(my4, opp4)

		rankMy4Final, _ := GetBadugiRank(my4Final)
		rankOpp4Final, _ := GetBadugiRank(opp4Final)

		if rankMy4Final < rankOpp4Final {
			potWins += 1.0
		} else if rankMy4Final == rankOpp4Final {
			potWins += 0.5
		}
	}

	equity := potWins / float64(iters)
	return equity
}

type BadugiRanking map[string]int

func GetBadugiRank(cards []Card) (int, error) {
	c, err := ExtractValidBadugiCards(cards)
	if err != nil {
		return 0, err
	}

	key, err := cardsToBadugiKey(c)
	if err != nil {
		return 0, err
	}

	ranking, _ := LoadBadugiRanking("badugi_ranking.json")
	rank, ok := ranking[key]
	if !ok {
		return 0, fmt.Errorf("hand %v (key: %s) not found in ranking", cards, key)
	}
	return rank, nil
}

// for GetBadugiRank() & changeForBadugi()
func ExtractValidBadugiCards(cards []Card) ([]Card, error) {
	if len(cards) != 4 {
		return nil, fmt.Errorf("need 4 cards, got %d", len(cards))
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank() < cards[j].Rank()
	})

	usedRanks := map[int]struct{}{}
	usedSuits := map[int]struct{}{}
	var valid []Card

	for _, c := range cards {
		r, s := c.Rank(), c.Suit()
		if _, ok := usedRanks[r]; ok {
			continue
		}
		if _, ok := usedSuits[s]; ok {
			continue
		}
		usedRanks[r] = struct{}{}
		usedSuits[s] = struct{}{}
		valid = append(valid, c)
	}

	return valid, nil
}

// for GetBadugiRank()
func cardsToBadugiKey(cards []Card) (string, error) {
	valid, err := ExtractValidBadugiCards(cards)
	if err != nil {
		return "", err
	}

	var parts []string
	for _, c := range valid {
		parts = append(parts, rankToChar[c.Rank()])
	}

	return strings.Join(parts, ""), nil
}

// for GetBadugiRank()
func LoadBadugiRanking(filename string) (BadugiRanking, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ranking BadugiRanking
	if err := json.Unmarshal(bytes, &ranking); err != nil {
		return nil, err
	}
	return ranking, nil
}

func changeForBadugi(my4 []Card, opp4 []Card) ([]Card, []Card) {
	if my4 != nil {
	}
	if opp4 != nil {
	}
	my4Final := []Card{0, 5, 10, 15}
	opp4Final := []Card{4, 9, 14, 19}
	return my4Final, opp4Final
}
