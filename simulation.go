package main

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

const (
	height = 7
	reps = 3
	concurrencyDepth = 2
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	g := NewGrid()
	g.PlaceRandom()
	g.PlaceRandom()

	fmt.Printf("%s\n", g)

	moves := 0

	for {
		avgBest := g.NextMove(height, reps, concurrencyDepth)
		if avgBest == -1 {
			fmt.Println("Couldn't find good direction. Game over")
			break
		}

		fmt.Printf("Moving %s\n", DirectionToString(avgBest))
		g.Move(avgBest)
		moves++

		fmt.Printf("%s\n", g)

		if !g.PlaceRandom() {
			fmt.Println("Couldn't place piece. Game over")
			break;
		}
	}

	fmt.Printf("After %d moves:\n", moves)
	fmt.Println(g)
}
