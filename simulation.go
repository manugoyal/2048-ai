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
	fmt.Println()

	const (
		height = 6
		reps = 3
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

		fmt.Println(g)
		fmt.Println()

		if !g.PlaceRandom() {
			fmt.Println("Couldn't place piece. Game over")
			break;
		}
		// if g.Score > 25000 {
		// 	fmt.Println("Reached maximum score. Game over")
		// 	break;
		// }
	}

	fmt.Printf("After %d moves:\n", moves)
	fmt.Println(g)
}
