package poker

import (
	"testing"
)

func TestSimulateEquity(t *testing.T) {
	tests := []struct {
		name        string
		game        Game
		hand        []Card
		minEquity   float64
		maxEquity   float64
		description string
	}{
		{
			name: "Four aces in Drawmaha-Hi",
			game: DrawmahaHi{},
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"),
			},
			minEquity:   0.99,
			maxEquity:   1.0,
			description: "Should have near 100% equity",
		},
		{
			name: "Low straight draw in Drawmaha-Hi",
			game: DrawmahaHi{},
			hand: []Card{
				mustCard("2s"), mustCard("3d"), mustCard("4h"), mustCard("5c"),
			},
			minEquity:   0.2,
			maxEquity:   0.4,
			description: "Should have low equity",
		},
		{
			name: "Perfect badugi hand",
			game: BadugiGame{},
			hand: []Card{
				mustCard("As"), mustCard("2d"), mustCard("3h"), mustCard("4c"),
			},
			minEquity:   0.95,
			maxEquity:   1.0,
			description: "A234 rainbow should dominate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equity := SimulateEquity(tt.game, tt.hand, 1000)
			if equity < tt.minEquity || equity > tt.maxEquity {
				t.Errorf("%s: equity = %.3f, want between %.3f and %.3f",
					tt.description, equity, tt.minEquity, tt.maxEquity)
			}
		})
	}
}

func TestPickBestGame(t *testing.T) {
	tests := []struct {
		name     string
		hand     []Card
		wantGame string
	}{
		{
			name: "Four aces should pick HiDuGi",
			hand: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"),
			},
			wantGame: "HiDuGi",
		},
		{
			name: "Low rainbow cards should pick Badugi",
			hand: []Card{
				mustCard("2s"), mustCard("3d"), mustCard("4h"), mustCard("5c"),
			},
			wantGame: "Badugi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			best, equities := PickBestGame(tt.hand, 1000)
			if best.Name() != tt.wantGame {
				t.Errorf("PickBestGame() selected %s, want %s", best.Name(), tt.wantGame)
				// Print all equities for debugging
				for game, eq := range equities {
					t.Logf("%s: %.3f", game, eq)
				}
			}
		})
	}
}
