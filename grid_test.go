package main

import "testing"

func TestClone(t *testing.T) {
	g := NewGrid()
	g.Tiles[0][0] = 2
	g.Score = 10
	g2 := g.Clone()

	if !g.Equal(g2) {
		t.Fatalf("clone doesn't produce same data")
	}

	g2.Tiles[0][1] = 20
	if g.Equal(g2) {
		t.Fatalf("clones aren't independent")
	}
}

func TestMovement(t *testing.T) {
	g := NewGrid()
	g.Tiles[0][0] = 2;
	g.Tiles[0][1] = 2;
	if !g.Equal(FromSparseGrid([]int{0, 0, 2, 0, 1, 2}, 0)) {
		t.Fatalf("initialization failed")
	}

	g.Move(RIGHT)
	if !g.Equal(FromSparseGrid([]int{0, 3, 4}, 4)) {
		t.Fatalf("move right 1 failed")
	}

	g.Tiles[1][2] = 4;
	g.Move(UP)
	if !g.Equal(FromSparseGrid([]int{0, 2, 4, 0, 3, 4}, 4)) {
		t.Fatalf("move up 1 failed")
	}

	g.Move(LEFT);
	if !g.Equal(FromSparseGrid([]int{0, 0, 8}, 12)) {
		t.Fatalf("move left 1 failed")
	}

	g.Tiles[1][0] = 2;
	g.Tiles[2][2] = 16;
	g.Tiles[3][2] = 16;
	g.Move(DOWN);
	if !g.Equal(FromSparseGrid([]int{3, 0, 2, 2, 0, 8, 3, 2, 32}, 44)) {
		t.Fatalf("move down 1 failed")
	}
}

func TestPlaceRandom(t *testing.T) {
	g := NewGrid()

	if !g.PlaceRandom() {
		t.Fatalf("PlaceRandom on empty grid should return true")
	}
	var total GridNum = 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			total += g.Tiles[r][c]
		}
	}
	if total != 2 {
		t.Fatalf("PlaceRandom didn't place exactly one 2")
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g.Tiles[r][c] = 4
		}
	}
	if g.PlaceRandom() {
		t.Fatalf("PlaceRandom on full grid should return false")
	}
}
