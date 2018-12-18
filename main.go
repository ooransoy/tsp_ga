package main

import (
	"flag"
	"fmt"
	"ga"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

//------------------------------------------------------------------------------

type Sequence [][2]float64

type LengthError struct{}

func (err LengthError) Error() string {
	return "Sequence lengths are not equal"
}

func init() {
	rand.Seed(time.Now().UnixNano())
	ga.Crossover = func(a []interface{}, b []interface{}) ([]interface{}, error) {
		if len(a) != len(b) {
			return []interface{}{}, LengthError{}
		}

		// http://www.rubicite.com/Tutorials/GeneticAlgorithms/CrossoverOperators/Order1CrossoverOperator.aspx
		p1 := rand.Intn(len(a) + 1) // Split point 1
		p2 := rand.Intn(len(a) + 1) // Split point 2

		if p1 > p2 {
			p1, p2 = p2, p1 // Ensure that p1 !> p2
		}

		// ----------------------- TODO optimize -----------------------
		o := make([]interface{}, len(a)) // Offspring

		// Copy a[p1,p2] to o[p1,p2]
		for i, g := range a[p1:p2] {
			o[i+p1] = g
		}

		// TODO explain spaghetti
		var oi int
		for i := 0; i < len(b); i++ {
			if oi == len(b) {
				break
			}
			if o[oi] != nil {
				oi++
				i--
				continue
			}
			for _, cg := range a[p1:p2] {
				if cg == b[i] {
					goto c
				}
			}
			o[oi] = b[i]
			oi++
		c:
			continue
		}
		// -----------/----------- TODO optimize -----------/-----------

		return o, nil
	}

	ga.Mutate = func(g []interface{}, _ []interface{}) ([]interface{}, error) {
		p := rand.Intn(len(g) - 1)
		g[p], g[p+1] = g[p+1], g[p]
		return g, nil
	}
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	//---------
	popSize := 500
	mutRate := 0.01
	selectK := 400
	evCount := 500

	solve(seq1, popSize, selectK, evCount, mutRate)
}

func solve(seq Sequence, popSize, selectK, evCount int, mutRate float64) {
	start := time.Now()
	// Create Population
	pop := make([][]interface{}, popSize)
	for i := range pop {
		pop[i] = shuffle(seq)
	}

	ff := func(s []interface{}) float64 {
		var sum float64
		var prev = s[len(s)-1].([2]float64)

		for _, c := range s {
			sum += dist(c.([2]float64), prev)
			prev = c.([2]float64)
		}

		return 1 / sum
	}

	for i := 0; i < evCount; i++ {
		fmt.Printf("Step %d/%d\n", i+1, evCount)
		var err error
		pop, err = ga.EvolvePop(pop, ff, selectK, []interface{}{}, mutRate)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\x1b[1F")
	}

	best := ga.Select(pop, ff)
	fmt.Println(1/ff(best), best)

	d := time.Since(start)
	fmt.Println("Solution benchmark:", d)
}

func shuffle(s Sequence) []interface{} {
	out := make([]interface{}, len(s))
	for i, n := range rand.Perm(len(s)) {
		out[i] = s[n]
	}

	return out
}

func dist(a, b [2]float64) float64 {
	return math.Pow(math.Pow(a[0]-b[0], 2)+math.Pow(a[1]-b[1], 2), 0.5)
}

// Sequence samples
var seq1 = Sequence{
	{0, 0},
	{1.1, 2.5},
	{-2.25, 2},
	{1, -1.2},
	{2.65, -1.75},
	{-1.5, -0.75},
	{3.4, 1.5},
	{-4, -3},
	{0.1, -2.6},
}
