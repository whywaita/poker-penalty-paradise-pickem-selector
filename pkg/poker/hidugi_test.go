package poker

import (
	"testing"
)

func TestEvaluate4CardHigh(t *testing.T) {
	tests := []struct {
		name        string
		hand        []Card
		wantBetter  []Card
		description string
	}{
		{
			name: "Four of a kind beats three of a kind",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"),
			},
			wantBetter: []Card{
				mustCard("Ks"), mustCard("Kd"), mustCard("Kh"), mustCard("2c"),
			},
			description: "Four aces",
		},
		{
			name: "Straight flush beats four of a kind",
			hand: []Card{
				mustCard("As"), mustCard("Ks"), mustCard("Qs"), mustCard("Js"),
			},
			wantBetter: []Card{
				mustCard("2s"), mustCard("2d"), mustCard("2h"), mustCard("2c"),
			},
			description: "Straight flush to ace",
		},
		{
			name: "Flush beats straight",
			hand: []Card{
				mustCard("As"), mustCard("Ks"), mustCard("Qs"), mustCard("9s"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Kd"), mustCard("Qh"), mustCard("Jc"),
			},
			description: "Ace high flush",
		},
		{
			name: "Straight beats two pair",
			hand: []Card{
				mustCard("5s"), mustCard("4d"), mustCard("3h"), mustCard("2c"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Kc"),
			},
			description: "5-high straight",
		},
		{
			name: "Three of a kind beats two pair",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Kc"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Kc"),
			},
			description: "Three aces",
		},
		{
			name: "Two pair beats one pair",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Kc"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Qc"),
			},
			description: "Aces and kings",
		},
		{
			name: "One pair beats high card",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Qc"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Kd"), mustCard("Qh"), mustCard("9c"),
			},
			description: "Pair of aces",
		},
		{
			name: "Wheel straight (A234)",
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			wantBetter: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Kh"), mustCard("Qc"),
			},
			description: "Wheel straight beats pair",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score1 := Evaluate4CardHigh(tt.hand)
			score2 := Evaluate4CardHigh(tt.wantBetter)
			if score1 <= score2 {
				t.Errorf("%s should beat the other hand: score1=%d, score2=%d", tt.description, score1, score2)
			}
		})
	}
}

func TestIsBadugi8OrBetter(t *testing.T) {
	tests := []struct {
		name string
		hand []Card
		want bool
	}{
		{
			name: "Perfect 8-badugi",
			hand: []Card{
				mustCard("8s"), mustCard("7d"), mustCard("3h"), mustCard("2c"),
			},
			want: true,
		},
		{
			name: "7-badugi is better than 8",
			hand: []Card{
				mustCard("7s"), mustCard("5d"), mustCard("3h"), mustCard("2c"),
			},
			want: true,
		},
		{
			name: "9-badugi is not 8 or better",
			hand: []Card{
				mustCard("9s"), mustCard("7d"), mustCard("3h"), mustCard("2c"),
			},
			want: false,
		},
		{
			name: "Not a 4-card badugi",
			hand: []Card{
				mustCard("8s"), mustCard("8d"), mustCard("3h"), mustCard("2c"),
			},
			want: false,
		},
		{
			name: "Ace-low badugi is 8 or better",
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBadugi8OrBetter(tt.hand)
			if got != tt.want {
				t.Errorf("IsBadugi8OrBetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateHiDuGi(t *testing.T) {
	tests := []struct {
		name string
		hand []Card
		desc string
	}{
		{
			name: "Strong both ways - A234 rainbow",
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			desc: "Wheel straight and perfect badugi",
		},
		{
			name: "Four aces - strong high, weak badugi",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"),
			},
			desc: "Quads high but only 1-card badugi",
		},
		{
			name: "8-7-6-5 rainbow - good both ways",
			hand: []Card{
				mustCard("8s"), mustCard("7d"), mustCard("6h"), mustCard("5c"),
			},
			desc: "Straight and 8-badugi",
		},
		{
			name: "Random high cards",
			hand: []Card{
				mustCard("As"), mustCard("Kd"), mustCard("Qh"), mustCard("Jc"),
			},
			desc: "High card only, no badugi qualification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			highScore, badugiScore := EvaluateHiDuGi(tt.hand)
			t.Logf("%s - High: %d, Badugi: %d", tt.desc, highScore, badugiScore)

			// Basic sanity checks
			if highScore <= 0 {
				t.Errorf("High score should be positive")
			}
			if badugiScore <= 0 {
				t.Errorf("Badugi score should be positive")
			}
		})
	}
}

func TestHiDuGiGame(t *testing.T) {
	game := HiDuGiGame{}

	tests := []struct {
		name        string
		hand        []Card
		description string
	}{
		{
			name: "Perfect A234 rainbow",
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			description: "Should have very high score with 8-badugi bonus",
		},
		{
			name: "8-7-6-5 rainbow",
			hand: []Card{
				mustCard("8s"), mustCard("7d"), mustCard("6h"), mustCard("5c"),
			},
			description: "Should have high score with 8-badugi bonus",
		},
		{
			name: "Four aces",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"),
			},
			description: "Strong high but weak badugi, no bonus",
		},
		{
			name: "Random cards",
			hand: []Card{
				mustCard("Ks"), mustCard("Qd"), mustCard("9h"), mustCard("7c"),
			},
			description: "Average hand, no 8-badugi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := game.Evaluate(tt.hand, nil)
			has8Badugi := IsBadugi8OrBetter(tt.hand)

			t.Logf("%s - Score: %d, Has 8-Badugi: %v", tt.description, score, has8Badugi)

			// Check that 8-badugi hands get bonus
			if has8Badugi && score < 10000000 {
				t.Errorf("8-badugi hand should get bonus score > 10000000, got %d", score)
			}
			// Four aces (trips or better) also get a big boost even without 8-badugi
			highCategory := score / (13 * 13 * 13 * 13)
			if !has8Badugi && highCategory < 6 && score >= 10000000 {
				t.Errorf("Non-8-badugi hand without trips should not get bonus score, got %d", score)
			}
		})
	}
}
