package main

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

var (
	threadNum int = runtime.NumCPU()
	height = 6
	reps = threadNum
)

func main() {
	runtime.GOMAXPROCS(threadNum)
	rand.Seed(time.Now().UnixNano())

	g := NewGrid()
	g.PlaceRandom()
	g.PlaceRandom()

	fmt.Printf("%s\n", g)

	moves := 0

	for {
		avgBest := g.NextMove(height, reps, threadNum)
		if avgBest == -1 {
			fmt.Println("Couldn't find good direction. Game over")
			break
		}

		fmt.Printf("Moving %s\n", DirectionToString(avgBest))
		g.Move(avgBest)
		moves++

		if !g.PlaceRandom() {
			fmt.Println("Couldn't place piece. Game over")
			break;
		}

		fmt.Printf("%s\n", g)

	}

	fmt.Printf("After %d moves:\n", moves)
	fmt.Println(g)
}
