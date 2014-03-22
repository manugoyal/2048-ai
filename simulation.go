package main;

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	g := NewGrid()
	g.PlaceRandom()
	g.PlaceRandom()

	fmt.Println(g)

	const (
		height = 7
		reps = 5
		concurrencyDepth = 2
	)

	moves := 0

	for {
		counts := map[int]int{LEFT: 0, RIGHT: 0, UP: 0, DOWN: 0}
		for j := 0; j < reps; j++ {
			t := NewTree()
			t.G = g.Clone()
			t.Fill(height, concurrencyDepth, nil)
			t.Score()
			counts[t.BestDirection]++
		}

		avgBest := -1
		avgOcc := 0
		for k, v := range counts {
			if v > avgOcc {
				avgBest = k
				avgOcc = v
			}
		}

		if avgBest == -1 {
			fmt.Printf("Couldn't find good direction, quitting\n")
			return
		}

		fmt.Printf("Moving %s\n", DirectionToString(avgBest))
		g.Move(avgBest)
		moves++

		if !g.PlaceRandom() {
			fmt.Println("Couldn't place piece. Game over")
			break;
		}

		fmt.Println(g)
	}

	fmt.Printf("After %d moves:\n", moves)
	fmt.Println(g)
}
