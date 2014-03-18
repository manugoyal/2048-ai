// The grid representation

package main

import (
	"fmt"
	"math/rand"
)

// Grid definition

// The number of rows and columns in the grid
const (
	rows = 4
	cols = 4
)

type GridNum uint16
type Grid struct {
	Tiles [][]GridNum
	Score uint32
}

// NewGrid returns an empty grid
func NewGrid() *Grid {
	g := make([][]GridNum, rows)
	for i := range g {
		g[i] = make([]GridNum, cols)
	}
	return &Grid{g, 0}
}

// FromSparseGrid takes a slice consisting of triples of integers,
// where the first is the row, the second the column, and the third
// the value, and creates a grid.
func FromSparseGrid(vals []int, score uint32) *Grid {
	g := NewGrid()
	g.Score = score
	for i := 0; i < len(vals); i += 3 {
		g.Tiles[vals[i]][vals[i+1]] = GridNum(vals[i+2])
	}
	return g
}

// Clone returns an independent copy of the grid
func (grid *Grid) Clone() *Grid {
	g := NewGrid()
	for i := range g.Tiles {
		copy(g.Tiles[i], grid.Tiles[i])
	}
	g.Score = grid.Score
	return g
}

// Places a 2 or 4 tile (90% chance it's a 2) at a random place in the
// board. Assumes the random number generator is already seeded. If
// there is no place for a tile, it returns false.
func (grid *Grid) PlaceRandom() bool {
	// The random number is a position to start searching at, which
	// wraps around.
	const total = rows * cols
	var tileval GridNum = 2
	if (rand.Float32() < 0.1) {
		tileval = 4;
	}
	startPos := rand.Intn(total)
	for i := (startPos+1) % total; i != startPos; i = (i + 1) % total {
		r, c := i / cols, i % cols
		if grid.Tiles[r][c] == 0 {
			grid.Tiles[r][c] = tileval
			return true
		}
	}
	r, c := startPos / cols, startPos % cols
	if grid.Tiles[r][c] == 0 {
		grid.Tiles[r][c] = tileval
		return true
	}
	return false
}

// String turns the grid into a printable representation
func (grid *Grid) String() string {
	ret := ""
	for i := range grid.Tiles {
		ret += fmt.Sprintf("%v\n", grid.Tiles[i])
	}
	return ret + "Score: " + fmt.Sprint(grid.Score)
}

// Equal returns true if both grids contain the same data, and false
// otherwise
func (grid *Grid) Equal(grid2 *Grid) bool {
	if grid.Score != grid2.Score {
		return false
	}
	for i := range grid.Tiles {
		for j := range grid.Tiles {
			if grid.Tiles[i][j] != grid2.Tiles[i][j] {
				return false
			}
		}
	}
	return true
}

// Movement directions
const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

// DirectionToString turns a direction into a string
func DirectionToString(direction int) string {
	switch direction {
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	default:
		return "Invalid direction"
	}
}

// Move moves the tiles in the grid in the given direction. It returns
// false if no tiles moved, and true otherwise.
func (grid *Grid) Move(direction int) (ret bool) {
	switch direction {
	case LEFT, RIGHT:
		// Goes row by row
		for r := 0; r < rows; r++ {
			// The column to merge with, start column, end column, and
			// column increment
			var moveCol, start, end, inc int
			if direction == LEFT {
				moveCol, start, end, inc = 0, 1, cols, 1
			} else {
				moveCol, start, end, inc = cols-1, cols-2, -1, -1
			}
			for c := start; c != end; c += inc {
				switch {
				case grid.Tiles[r][c] == 0:
					continue
				case grid.Tiles[r][moveCol] == 0:
					// Move grid.Tiles[r][c] all the way down
					grid.Tiles[r][c], grid.Tiles[r][moveCol] = 0, grid.Tiles[r][c]
					ret = true
				case grid.Tiles[r][moveCol] == grid.Tiles[r][c]:
					// Merge grid.Tiles[r][c] with grid.Tiles[r][moveCol]
					grid.Tiles[r][c], grid.Tiles[r][moveCol] = 0, grid.Tiles[r][c]*2
					grid.Score += uint32(grid.Tiles[r][moveCol])
					ret = true
				default:
					// Increment moveCol and move grid.Tiles[r][c] there, if
					// it isn't already
					moveCol += inc
					if moveCol != c {
						grid.Tiles[r][c], grid.Tiles[r][moveCol] = grid.Tiles[r][moveCol], grid.Tiles[r][c]
						ret = true
					}
				}
			}
		}
	case UP, DOWN:
		// Goes column by column
		for c := 0; c < cols; c++ {
			// The row to merge with, start row, end row, and row
			// increment
			var moveRow, start, end, inc int
			if direction == UP {
				moveRow, start, end, inc = 0, 1, rows, 1
			} else {
				moveRow, start, end, inc = rows-1, rows-2, -1, -1
			}
			for r := start; r != end; r += inc {
				switch {
				case grid.Tiles[r][c] == 0:
					continue
				case grid.Tiles[moveRow][c] == 0:
					// Move grid.Tiles[r][c] all the way down
					grid.Tiles[r][c], grid.Tiles[moveRow][c] = 0, grid.Tiles[r][c]
					ret = true
				case grid.Tiles[moveRow][c] == grid.Tiles[r][c]:
					// Merge grid.Tiles[r][c] with grid.Tiles[r][moveCol]
					grid.Tiles[r][c], grid.Tiles[moveRow][c] = 0, grid.Tiles[r][c]*2
					grid.Score += uint32(grid.Tiles[moveRow][c])
					ret = true
				default:
					// Increment moveRow and move grid.Tiles[r][c] there, if
					// it isn't already
					moveRow += inc
					if moveRow != r {
						grid.Tiles[r][c], grid.Tiles[moveRow][c] = grid.Tiles[moveRow][c], grid.Tiles[r][c]
						ret = true
					}
				}
			}
		}
	default:
		panic("Invalid direction")
	}
	return
}
