package poker

import (
	"testing"
)

func TestEvaluate5CardHigh(t *testing.T) {
	tests := []struct {
		name        string
		hand        []Card
		wantBetter  []Card // hand that should be worse
		description string
	}{
		{
			name: "Royal flush beats straight flush",
			hand: []Card{
				mustCard("As"), mustCard("Ks"), mustCard("Qs"), mustCard("Js"), mustCard("Ts"),
			},
			wantBetter: []Card{
				mustCard("9s"), mustCard("8s"), mustCard("7s"), mustCard("6s"), mustCard("5s"),
			},
			description: "Royal flush",
		},
		{
			name: "Four of a kind beats full house",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"), mustCard("2s"),
			},
			wantBetter: []Card{
				mustCard("Ks"), mustCard("Kd"), mustCard("Kh"), mustCard("Ac"), mustCard("As"),
			},
			description: "Four aces",
		},
		{
			name: "Full house beats flush",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("2c"), mustCard("2s"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ks"), mustCard("Qs"), mustCard("Js"), mustCard("8s"),
			},
			description: "Aces full of twos",
		},
		{
			name: "Flush beats straight",
			hand: []Card{
				mustCard("As"), mustCard("Ks"), mustCard("Qs"), mustCard("Js"), mustCard("9s"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Kd"), mustCard("Qh"), mustCard("Jc"), mustCard("Ts"),
			},
			description: "Ace high flush",
		},
		{
			name: "Straight beats three of a kind",
			hand: []Card{
				mustCard("5s"), mustCard("4d"), mustCard("3h"), mustCard("2c"), mustCard("As"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Kc"), mustCard("Qs"),
			},
			description: "Wheel straight",
		},
		{
			name: "Three of a kind beats two pair",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Kc"), mustCard("Qs"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Kc"), mustCard("Qs"),
			},
			description: "Three aces",
		},
		{
			name: "Two pair beats one pair",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Kc"), mustCard("Qs"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Qc"), mustCard("Js"),
			},
			description: "Aces and kings",
		},
		{
			name: "One pair beats high card",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Qc"), mustCard("Js"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Kd"), mustCard("Qh"), mustCard("Jc"), mustCard("9s"),
			},
			description: "Pair of aces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score1 := Evaluate5CardHigh(tt.hand)
			score2 := Evaluate5CardHigh(tt.wantBetter)
			if score1 <= score2 {
				t.Errorf("%s should beat the other hand: score1=%d, score2=%d", tt.description, score1, score2)
			}
		})
	}
}

func TestEvaluateBadugi(t *testing.T) {
	tests := []struct {
		name        string
		hand        []Card
		wantBetter  []Card
		description string
	}{
		{
			name: "4-card badugi beats 3-card badugi",
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("3c"),
			},
			description: "A234 rainbow",
		},
		{
			name: "Lower 4-card badugi beats higher",
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("5c"),
			},
			description: "A234 beats A235",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score1 := EvaluateBadugi(tt.hand)
			score2 := EvaluateBadugi(tt.wantBetter)
			if score1 <= score2 {
				t.Errorf("%s should beat the other hand: score1=%d, score2=%d", tt.description, score1, score2)
			}
		})
	}
}

// Helper function for tests
func mustCard(s string) Card {
	c, err := CardFromString(s)
	if err != nil {
		panic(err)
	}
	return c
}
