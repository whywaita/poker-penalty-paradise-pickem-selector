package main

import (
	"fmt"
	"os"
	"time"

	"github.com/whywaita/poker-penalty-paradise-pickem-selector/pkg/poker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s \"Ac Kd 2h 3c\"\n", os.Args[0])
		os.Exit(1)
	}
	hand, err := poker.ParseHand(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	start := time.Now()
	best, eqs := poker.PickBestGame(hand, 20000) // 20k sims per game ≈200‑300 ms
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
