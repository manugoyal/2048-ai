// Creates a game tree and tries to make decisions that maximize the
// score. It's not a minimax algorithm, since the opponent's decisions
// are random.
package main

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Tree struct {
	G             *Grid
	Children      []*Tree
	BestScore     uint32
	BestDirection int
}

// NewTree returns an tree with an empty grid
func NewTree(g *Grid) *Tree {
	return &Tree{g, nil, g.Score, -1}
}

// Given the height of the tree, it will fill out the tree to nodes of
// height 0.
func (t *Tree) Fill(height int, localRand *rand.Rand) {
	if height == 0 {
		return
	}
	// Generate new children
	t.Children = make([]*Tree, 4)

	t.BestScore = 0
	for i := 0; i < 4; i++ {
		node := NewTree(t.G.Clone())
		// We only execute the move if tiles would be moving
		if node.G.Move(i) {
			t.Children[i] = node
			if node.G.PlaceRandom(localRand) {
				node.Fill(height-1, localRand)
			}
			if node.BestScore > t.BestScore {
				t.BestScore, t.BestDirection = node.BestScore, i
			}
		}
	}
}

// Given a grid and some parameters, it figures out the next best
// move. If it returns -1, that means it couldn't find a move.
func (g *Grid) NextMove(height, reps, threadNum int) int {
	baseReps, leftoverReps := reps/threadNum, reps%threadNum
	iterations := threadNum
	if baseReps == 0 {
		iterations, leftoverReps = leftoverReps, 0
	}
	// directionCounts keeps track of how many reps reported each
	// direction as the best
	directionCounts := make([]int64, 4)
	var wg sync.WaitGroup
	t := time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		// This goroutine creates a tree repsPerThread times and
		// adds each resulting direction to the directions channel
		wg.Add(1)
		go func(i int) {
			repsPerThread := baseReps
			if i < leftoverReps {
				repsPerThread++
			}
			localRand := rand.New(rand.NewSource(t + int64(i)))
			for j := 0; j < repsPerThread; j++ {
				t := NewTree(g.Clone())
				t.Fill(height, localRand)
				if t.BestDirection >= 0 && t.BestDirection < 4 {
					atomic.AddInt64(&directionCounts[t.BestDirection], 1)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	bestDirection := -1
	var maxOcc int64
	for direction, occ := range directionCounts {
		if occ > maxOcc {
			bestDirection, maxOcc = direction, occ
		}
	}
	return bestDirection
}
