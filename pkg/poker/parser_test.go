package poker

import (
	"reflect"
	"testing"
)

func TestParseHand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Card
		wantErr bool
	}{
		{
			name:  "Space separated",
			input: "As Ad Ah Ac",
			want: []Card{
				mustCard("As"), mustCard("Ad"), mustCard("Ah"), mustCard("Ac"),
			},
			wantErr: false,
		},
		{
			name:  "Comma separated",
			input: "2s,3d,4h,5c",
			want: []Card{
				mustCard("2s"), mustCard("3d"), mustCard("4h"), mustCard("5c"),
			},
			wantErr: false,
		},
		{
			name:    "Too few cards",
			input:   "As Ad Ah",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Too many cards",
			input:   "As Ad Ah Ac Ks",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Duplicate card",
			input:   "As As Ah Ac",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid card",
			input:   "As Ad Ah Xx",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseHand() = %v, want %v", got, tt.want)
			}
		})
	}
}
