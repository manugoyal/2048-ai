// Creates a game tree and tries to make decisions that maximize the
// score. It's not a minimax algorithm, since the opponent's decisions
// are random.
package main

import "math"

type Tree struct {
	G *Grid
	Children []*Tree
	Done bool
	BestScore uint32
	BestDirection int
}

// NewTree returns an tree with a starting grid, that is, a grid with
// two tiles of value 2 placed randomly.
func NewTree() *Tree {
	ret := Tree{NewGrid(), nil, false, 0, 0}
	ret.G.PlaceRandom()
	ret.G.PlaceRandom()
	return &ret
}

func generateFinChan(concurrencyDepth int) chan bool {
	if concurrencyDepth > 0 {
		return make(chan bool, 4)
	}
	return nil
}

// Given the height of the tree, it will fill out the tree to nodes of
// height 0. If the tree already has children, it won't generate new
// ones, but it will recursively call Fill. This should allow for
// iterative deepening.
func (t *Tree) Fill(height, concurrencyDepth int, fin chan bool) {
	if height == 0 {
		return
	}
	if t.Children == nil {
		// Generate new children. If it's within the concurrency
		// depth, do it concurrently
		subfin := generateFinChan(concurrencyDepth)
		fills := 0

		t.Children = make([]*Tree, 4)
		
		for i := 0; i < 4; i++ {
			node := &Tree{t.G.Clone(), nil, false, 0, 0}
			node.G.Move(i)
			t.Children[i] = node

			if !node.G.PlaceRandom() {
				node.Done = true
			} else {
				if concurrencyDepth > 0 {
					go node.Fill(height-1, concurrencyDepth-1, subfin)
					fills++
				} else {
					node.Fill(height-1, concurrencyDepth, subfin)
				}
			}
		}
		
		for i := 0; i < fills; i++ {
			<-subfin
		}
	} else {
		// Recursively fill the children that aren't done
		for _, node := range t.Children {
			if !node.Done {
				node.Fill(height-1, concurrencyDepth-1, generateFinChan(concurrencyDepth))
			}
		}
	}
	if fin != nil {
		fin <- true
	}
}

// Given a filled tree, it takes the scores of the leaf nodes and
// fills up the scores and directions of the parent nodes using the
// minimax algorithm. The root of the tree will contain the best score
// and direction
func (t *Tree) Score() {
	if t.Children == nil || t.Done {
		t.BestScore = t.G.Score
		t.BestDirection = -1
	} else {
		t.BestScore = math.MaxUint32
		t.BestDirection = -1
		for i, child := range t.Children {
			child.Score()
			if child.BestScore < t.BestScore {
				t.BestScore = child.BestScore
				t.BestDirection = i
			}
		}
	}
}










