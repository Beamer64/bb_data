// Command cli prints a random joke. Useful as a smoke test that the
// embedded datasets compiled in correctly: `go run ./examples/cli`.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Beamer64/bb_data"
	"github.com/Beamer64/bb_data/affirmations"
	"github.com/Beamer64/bb_data/jokes"
)

func main() {
	start := time.Now()
	if err := bb_data.Load(); err != nil {
		log.Fatalf("bb_data.Load: %v", err)
	}
	fmt.Printf("loaded in %s\n", time.Since(start))
	fmt.Println("--- random joke ---")
	fmt.Println(jokes.Random())
	fmt.Println("--- random affirmation ---")
	fmt.Println(affirmations.Random())
}
