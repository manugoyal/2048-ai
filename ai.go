// Creates a game tree and tries to make decisions that maximize the
// score. It's not a minimax algorithm, since the opponent's decisions
// are random.
package main

import (
	"math/rand"
	"time"
)

type Tree struct {
	G             *Grid
	Children      []*Tree
	BestScore     uint32
	BestDirection int
}

// NewTree returns an tree with an empty grid
func NewTree() *Tree {
	return &Tree{NewGrid(), nil, 0, 0}
}

// Given the height of the tree, it will fill out the tree to nodes of
// height 0.
func (t *Tree) Fill(height int, localRand *rand.Rand) {
	if height == 0 {
		return
	}
	// Generate new children
	t.Children = make([]*Tree, 4)

	for i := 0; i < 4; i++ {
		node := &Tree{t.G.Clone(), nil, 0, 0}
		if node.G.Move(i) {
			// We only execute the move if tiles would be moving
			t.Children[i] = node

			if node.G.PlaceRandom(localRand) {
				node.Fill(height-1, localRand)
			}
		}
	}
}

// Given a filled tree, it takes the scores of the leaf nodes and
// fills up the scores and directions of the parent nodes using the
// minimax algorithm. The root of the tree will contain the best score
// and direction
func (t *Tree) Score() {
	if t == nil {
		return
	}
	if t.Children == nil {
		t.BestScore = t.G.Score
		t.BestDirection = -1
	} else {
		t.BestScore = 0
		t.BestDirection = -1
		for i, child := range t.Children {
			if child != nil {
				child.Score()
				if child.BestScore > t.BestScore {
					t.BestScore = child.BestScore
					t.BestDirection = i
				}
			}
		}
	}
}

// Given a grid and some parameters, it figures out the next best
// move. If it returns -1, that means it couldn't find a move.
func (g *Grid) NextMove(height, reps, threadNum int) int {
	directions := make(chan int)
	bestDirection := make(chan int)
	// This goroutine accumulates the directions from generating trees
	// and returns the best direction on the bestDirection channel.
	go func() {
		counts := map[int]int{LEFT: 0, RIGHT: 0, UP: 0, DOWN: 0}
		for dir := range directions {
			counts[dir]++
		}
		maxDir, maxOcc := 0, 0
		for direction, occurences := range counts {
			if occurences > maxOcc {
				maxDir, maxOcc = direction, occurences
			}
		}
		bestDirection <- maxDir
	}()
	baseReps, leftoverReps := reps/threadNum, reps%threadNum
	iterations := threadNum
	if baseReps == 0 {
		iterations = leftoverReps
	}
	loopChan := make(chan bool, iterations)
	t := time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		// This goroutine creates a tree repsPerThread times and
		// adds each resulting direction to the directions channel
		go func(i int) {
			repsPerThread := baseReps
			if i < leftoverReps {
				repsPerThread++
			}
			localRand := rand.New(rand.NewSource(t + int64(i)))
			for j := 0; j < repsPerThread; j++ {
				t := NewTree()
				t.G = g.Clone()
				t.Fill(height, localRand)
				t.Score()
				if t.BestDirection >= 0 && t.BestDirection < 4 {
					directions <- t.BestDirection
				}
			}
			loopChan <- true
		}(i)
	}
	for i := 0; i < iterations; i++ {
		<-loopChan
	}
	close(directions)
	return <-bestDirection
}
