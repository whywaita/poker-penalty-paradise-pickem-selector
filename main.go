// poker_selector.go
// Reference implementation: Given a 4‑card starting hand, estimate which mixed poker variant
// (Drawmaha / Badugi family / Omaha DoubleBoard, etc.) offers the highest equity versus a
// single random opponent by Monte‑Carlo simulation.
//
// ⚠️  This is a **prototype** meant to illustrate overall architecture in Go.
//     * Drawmaha‑Hi (5‑card high) and Badugi evaluators are fully functional.
//     * Other variants are stubbed—extend the `Game` interface implementations.
//     * The 5‑card evaluator is a simple hand‑category comparator (not the fastest LUT).
//
// Compile & run:
//     go run poker_selector.go "Ac Kd 2h 3c"

package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	mrand "math/rand"
	"os"
	"strings"
	"time"
)

// ===== Card / Deck helpers ==================================================

// Card is an int 0‑51 where rank = c/4 (0=2 .. 12=Ace) and suit = c%4 (0=♣,1=♦,2=♥,3=♠)
// This mapping makes bit‑operations convenient.

type Card int

var rankToChar = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
var suitToChar = []string{"c", "d", "h", "s"}

func cardFromString(s string) (Card, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid card string %q", s)
	}
	r, su := strings.ToUpper(string(s[0])), strings.ToLower(string(s[1]))
	var rank, suit int = -1, -1
	for i, c := range rankToChar {
		if strings.EqualFold(c, r) {
			rank = i
			break
		}
	}
	for i, c := range suitToChar {
		if c == su {
			suit = i
			break
		}
	}
	if rank == -1 || suit == -1 {
		return 0, fmt.Errorf("invalid card %q", s)
	}
	return Card(rank*4 + suit), nil
}

func (c Card) Rank() int { return int(c) / 4 }
func (c Card) Suit() int { return int(c) % 4 }

func (c Card) String() string {
	return rankToChar[c.Rank()] + suitToChar[c.Suit()]
}

func fullDeck() []Card {
	d := make([]Card, 52)
	for i := 0; i < 52; i++ {
		d[i] = Card(i)
	}
	return d
}

// removeCards returns deck with specified cards removed.
func removeCards(deck []Card, toRemove map[Card]struct{}) []Card {
	out := make([]Card, 0, len(deck))
	for _, c := range deck {
		if _, ok := toRemove[c]; !ok {
			out = append(out, c)
		}
	}
	return out
}

// randomIntn uses crypto/rand seeded math/rand for reproducibility.
var rng *mrand.Rand

func init() {
	var seed int64
	binary.Read(rand.Reader, binary.LittleEndian, &seed)
	rng = mrand.New(mrand.NewSource(seed))
}

// drawRandom draws n unique random cards from `deck` in‑place (Fisher‑Yates shuffle prefix).
func drawRandom(deck []Card, n int) ([]Card, []Card) {
	if n > len(deck) {
		panic("drawRandom: not enough cards")
	}
	for i := 0; i < n; i++ {
		j := rng.Intn(len(deck)-i) + i
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck[:n], deck[n:]
}

// ===== Hand evaluators =======================================================

// 5‑Card High evaluator – category ranking (higher is better).
// Category weights ensure unique ordering.
const (
	highCard = iota
	onePair
	twoPair
	trips
	straight
	flush
	fullHouse
	quads
	straightFlush
)

func evaluate5CardHigh(hand []Card) int64 {
	if len(hand) != 5 {
		panic("evaluate5CardHigh expects 5 cards")
	}
	ranks := make([]int, 13)
	suits := make([]int, 4)
	for _, c := range hand {
		ranks[c.Rank()]++
		suits[c.Suit()]++
	}
	isFlush := false
	for _, s := range suits {
		if s == 5 {
			isFlush = true
			break
		}
	}
	// detect straight (wheel A2345)
	top := -1
	consec := 0
	for r := 12; r >= 0; r-- {
		if ranks[r] > 0 {
			consec++
			if consec == 5 {
				top = r + 4 // highest rank of straight
				break
			}
		} else {
			consec = 0
		}
	}
	// wheel check
	if consec == 4 && ranks[12] > 0 && ranks[3] > 0 && ranks[2] > 0 && ranks[1] > 0 && ranks[0] > 0 {
		top = 3 // straight to 5
	}
	isStraight := top != -1

	// rank histogram counts
	counts := make(map[int]int) // rank -> count
	for r, cnt := range ranks {
		if cnt > 0 {
			counts[cnt] = counts[cnt]*13 + r // encode kicks by shifting
		}
	}

	var cat int64
	switch {
	case isStraight && isFlush:
		cat = straightFlush
	case counts[4] != 0:
		cat = quads
	case counts[3] != 0 && counts[2] != 0:
		cat = fullHouse
	case isFlush:
		cat = flush
	case isStraight:
		cat = straight
	case counts[3] != 0:
		cat = trips
	case counts[2] != 0 && counts[2]%13 != 0: // at least two pairs
		cat = twoPair
	case counts[2] != 0:
		cat = onePair
	default:
		cat = highCard
	}

	// naive kicker encoding: combine sorted ranks descending into value
	vals := make([]int, 0, 5)
	for r := 12; r >= 0; r-- {
		for i := 0; i < ranks[r]; i++ {
			vals = append(vals, r)
		}
	}
	var kicker int64
	for _, v := range vals {
		kicker = kicker*13 + int64(v)
	}
	return cat*int64(13*13*13*13*13) + kicker // category has highest weight
}

// Badugi evaluator – lower is better (4‑card low with unique suits).
// We convert to HIGH by negating within a large offset so that higher = better overall.
func evaluateBadugi(hand []Card) int64 {
	// select up to 4 cards with distinct suits and ranks minimal lexicographically.
	chosen := make([]Card, 0, 4)
	usedSuit := [4]bool{}
	usedRank := [13]bool{}
	// brute force: iterate ranks ascending, pick first available card with unique suit.
	for r := 0; r < 13 && len(chosen) < 4; r++ {
		for s := 0; s < 4 && len(chosen) < 4; s++ {
			c := Card(r*4 + s)
			for _, h := range hand {
				if h == c && !usedSuit[s] && !usedRank[r] {
					chosen = append(chosen, h)
					usedSuit[s] = true
					usedRank[r] = true
					break
				}
			}
		}
	}
	// score: fewer cards worse; then ranks high worse.
	score := int64(len(chosen)) * 1000
	for _, c := range chosen {
		score = score*13 + int64(c.Rank())
	}
	// invert so bigger is better (align with evaluate5CardHigh)
	return math.MaxInt64/2 - score
}

// ===== Game interface =======================================================

type Game interface {
	Name() string
	// CompleteHand fills in missing private and public cards for simulation.
	CompleteHand(my []Card, deck []Card) (myComplete []Card, oppHand []Card, board []Card, deckOut []Card)
	Evaluate(myComplete []Card, board []Card) int64
}

// ---------- Drawmaha Hi implementation -------------------------------------

type drawmahaHi struct{}

func (d drawmahaHi) Name() string { return "Drawmaha-Hi" }

func (d drawmahaHi) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	// Drawmaha deals 5‑card hands, no community board. Each player keeps 4 original cards & draws 1.
	var myHand, oppHand []Card
	var drawn []Card
	// Draw 1 card for hero
	drawn, deck = drawRandom(deck, 1)
	myHand = append(append([]Card(nil), my...), drawn...)
	// Opponent 5 cards
	oppHand, deck = drawRandom(deck, 5)
	return myHand, oppHand, nil, deck
}

func (d drawmahaHi) Evaluate(h []Card, board []Card) int64 {
	// Already completed to 5 cards
	return evaluate5CardHigh(h)
}

// ---------- Badugi implementation ------------------------------------------

type badugiGame struct{}

func (b badugiGame) Name() string { return "Badugi" }

func (b badugiGame) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	// Badugi uses 4‑card hands; hero already has 4.
	oppHand, deck := drawRandom(deck, 4)
	return my, oppHand, nil, deck
}

func (b badugiGame) Evaluate(h []Card, board []Card) int64 { return evaluateBadugi(h) }

// ---------- Stub games (extend as needed) -----------------------------------

type stubGame struct{ name string }

func (s stubGame) Name() string { return s.name }
func (s stubGame) CompleteHand(my []Card, deck []Card) ([]Card, []Card, []Card, []Card) {
	oppHand, deck := drawRandom(deck, len(my))
	return my, oppHand, nil, deck
}
func (s stubGame) Evaluate(h []Card, board []Card) int64 { return 0 }

// ===== Simulation & Selection ==============================================

// simulateEquity returns hero win probability against one opponent by Monte‑Carlo.
func simulateEquity(g Game, my4 []Card, iters int) float64 {
	wins, ties := 0, 0
	for i := 0; i < iters; i++ {
		deck := removeCards(fullDeck(), toSet(my4))
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

func toSet(cards []Card) map[Card]struct{} {
	m := make(map[Card]struct{}, len(cards))
	for _, c := range cards {
		m[c] = struct{}{}
	}
	return m
}

func pickBestGame(my4 []Card, iters int) (best Game, equities map[string]float64) {
	games := []Game{
		drawmahaHi{},
		badugiGame{},
		stubGame{"HiDuGi"},
		stubGame{"Drawmaha-2-7"},
		stubGame{"Prime"},
		stubGame{"Omaha DoubleBoard"},
	}
	equities = make(map[string]float64, len(games))
	for _, g := range games {
		eq := simulateEquity(g, my4, iters)
		equities[g.Name()] = eq
		if best == nil || eq > equities[best.Name()] {
			best = g
		}
	}
	return
}

// ===== CLI =================================================================

func parseHand(arg string) ([]Card, error) {
	parts := strings.Fields(arg)
	if len(parts) == 1 {
		// maybe comma‑separated "Ac,Kd,2h,3c"
		parts = strings.Split(arg, ",")
	}
	if len(parts) != 4 {
		return nil, fmt.Errorf("need exactly 4 cards, got %d", len(parts))
	}
	hand := make([]Card, 0, 4)
	seen := map[Card]struct{}{}
	for _, p := range parts {
		c, err := cardFromString(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		if _, dup := seen[c]; dup {
			return nil, fmt.Errorf("duplicate card %s", c)
		}
		seen[c] = struct{}{}
		hand = append(hand, c)
	}
	return hand, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s \"Ac Kd 2h 3c\"\n", os.Args[0])
		os.Exit(1)
	}
	hand, err := parseHand(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	start := time.Now()
	best, eqs := pickBestGame(hand, 20000) // 20k sims per game ≈200‑300 ms
	dur := time.Since(start)

	fmt.Printf("Hand: %s %s %s %s\n", hand[0], hand[1], hand[2], hand[3])
	fmt.Println("--------------------------------------------------")
	fmt.Println("Estimated equities vs 1 random opponent:")
	for g, e := range eqs {
		fmt.Printf("%-20s %.3f\n", g, e)
	}
	fmt.Println("--------------------------------------------------")
	fmt.Printf("=> Best game to register: %s\n", best.Name())
	fmt.Printf("Simulation time: %v\n", dur)
}
