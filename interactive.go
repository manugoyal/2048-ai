package main;

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
	"log"
)

func RequestTile() (row int, col int, val GridNum, err error) {
	fmt.Print("Enter the row, column, and value of the tile, separated by spaces: ")
	_, err = fmt.Scanf("%d %d %d", &row, &col, &val)
	if row < 0 || row >= rows {
		err = fmt.Errorf("row is out of bounds")
	}
	if col < 0 || col >= cols {
		err = fmt.Errorf("column is out of bounds")
	}
	return
}

func RequestDirection() (direction int, err error) {
	var dirchar string
	fmt.Print("Enter which direction you want to go (l, r, u, d): ")
	_, err = fmt.Scanf("%s", &dirchar)
	if err != nil {
		return
	}
	switch dirchar {
	case "l":
		direction = LEFT
	case "r":
		direction = RIGHT
	case "u":
		direction = UP
	case "d":
		direction = DOWN
	default:
		err = fmt.Errorf("Invalid direction")
	}
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	g := NewGrid()

	fmt.Println(g)
	fmt.Println("Requesting initial game state")
	row, col, val, err := RequestTile()
	if err != nil {
		log.Fatal(err)
	}
	g.Tiles[row][col] = val
	row, col, val, err = RequestTile()
	if err != nil {
		log.Fatal(err)
	}
	g.Tiles[row][col] = val

	fmt.Println(g)

	const (
		height = 9
		reps = 1
		concurrencyDepth = 2
	)

	defer func() {	fmt.Println(g) }()

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
			fmt.Printf("I couldn't find a good direction")
		} else {
			fmt.Printf("I recommend moving %s\n", DirectionToString(avgBest))
		}

		direction, err := RequestDirection()
		if err != nil {
			log.Fatal(err)
		}

		g.Move(direction)
		fmt.Println(g)

		row, col, val, err = RequestTile()
		if err != nil {
			log.Fatal(err)
		}
		if g.Tiles[row][col] != 0 {
			log.Fatal("Can't place tile there")
		} else {
			g.Tiles[row][col] = val
		}

		// if g.Score > 25000 {
		// 	fmt.Println("Reached maximum score. Game over")
		// 	break;
		// }
	}

}
