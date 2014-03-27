package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	threadNum int = runtime.NumCPU()
	height        = 6
	reps          = threadNum
)

func main() {
	runtime.GOMAXPROCS(threadNum)
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := NewGrid()
	g.PlaceRandom(localRand)
	g.PlaceRandom(localRand)

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

		if !g.PlaceRandom(localRand) {
			fmt.Println("Couldn't place piece. Game over")
			break
		}

		fmt.Printf("%s\n", g)

	}

	fmt.Printf("After %d moves:\n", moves)
	fmt.Println(g)
}
