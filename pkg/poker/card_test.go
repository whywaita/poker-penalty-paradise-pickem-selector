package poker

import (
	"testing"
)

func TestCardFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Card
		wantErr bool
	}{
		{"Ace of spades", "As", Card(12*4 + 3), false},
		{"Two of clubs", "2c", Card(0*4 + 0), false},
		{"King of hearts", "Kh", Card(11*4 + 2), false},
		{"Ten of diamonds", "Td", Card(8*4 + 1), false},
		{"Invalid card", "Xx", 0, true},
		{"Invalid suit", "Ax", 0, true},
		{"Invalid rank", "Xs", 0, true},
		{"Too short", "A", 0, true},
		{"Too long", "Asd", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CardFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CardFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardRankSuit(t *testing.T) {
	tests := []struct {
		card     Card
		wantRank int
		wantSuit int
	}{
		{Card(0), 0, 0},   // 2c
		{Card(51), 12, 3}, // As
		{Card(16), 4, 0},  // 6c
		{Card(45), 11, 1}, // Kd
	}

	for _, tt := range tests {
		t.Run(tt.card.String(), func(t *testing.T) {
			if got := tt.card.Rank(); got != tt.wantRank {
				t.Errorf("Card.Rank() = %v, want %v", got, tt.wantRank)
			}
			if got := tt.card.Suit(); got != tt.wantSuit {
				t.Errorf("Card.Suit() = %v, want %v", got, tt.wantSuit)
			}
		})
	}
}

func TestCardString(t *testing.T) {
	tests := []struct {
		card Card
		want string
	}{
		{Card(0), "2c"},
		{Card(51), "As"},
		{Card(12), "5c"},
		{Card(13), "5d"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.card.String(); got != tt.want {
				t.Errorf("Card.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
