package main

import (
	"flag"
	"fmt"
	"ga"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	gc "github.com/ooransoy/gocanvas"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

//------------------------------------------------------------------------------

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Sequence [][2]float64

type LengthError struct{}

func (err LengthError) Error() string {
	return "Sequence lengths are not equal"
}

func init() {
	ga.Crossover = func(a []interface{}, b []interface{}) ([]interface{}, error) {
		if len(a) != len(b) {
			return []interface{}{}, LengthError{}
		}

		// http://www.rubicite.com/Tutorials/GeneticAlgorithms/CrossoverOperators/Order1CrossoverOperatorandaspx
		p1 := r.Intn(len(a) + 1) // Split point 1
		p2 := r.Intn(len(a) + 1) // Split point 2

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
		p := r.Intn(len(g) - 1)
		g[p], g[p+1] = g[p+1], g[p]
		return g, nil
	}

	gc.Width = 800
	gc.Height = 600
	gc.Tick = time.Millisecond
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
	popSize := 800
	mutRate := 0.5
	selectK := 640
	evCount := 25000

	gc.Loop = func() {}

	go func(){
		for {
			solve(seq1, popSize, selectK, evCount, mutRate)
		}
	}()
	gc.Start()
}

var count int

func solve(seq Sequence, popSize, selectK, evCount int, mutRate float64) {
	start := time.Now()
	// Create Population
	pop := make([][]interface{}, popSize)
	for i := range pop {
		pop[i] = shuffle(seq)
	}

	ff := func(s []interface{}) float64 {
		var sum float64
		var prev = cast(s[len(s)-1])

		for _, c := range s {
			sum += dist(cast(c), prev)
			prev = cast(c)
		}

		return 1 / sum
	}

	cb := make([][2]float64, len(seq)) // Current best
	c := color.RGBA{uint8(r.Intn(256)), uint8(r.Intn(256)), uint8(r.Intn(256)), 255}
	gc.Loop = func() {
		gc.ResetCanvas()
		for i, g := range ga.Select(pop, ff) {
			cb[i] = cast(g)
		}
		gc.ScatterPlot(cb, c, [2]float64{75, -75}, gc.ClosedConnect)
	}

	for i := 0; i < evCount; i++ {
		//fmt.Printf("Step %d/%d\n", i+1, evCount)
		var err error
		pop, err = ga.EvolvePop(pop, ff, selectK, []interface{}{}, mutRate)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("\x1b[1F")
	}

	best := ga.Select(pop, ff)
	fmt.Println(1/ff(best), best)

	d := time.Since(start)
	fmt.Println("Solution benchmark:", d)
	count++
	fmt.Println("--------", count, "--------")
}

func shuffle(s Sequence) []interface{} {
	out := make([]interface{}, len(s))
	for i := range out {
		out[i] = s[i]
	}
	r.Shuffle(len(s), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

func dist(a, b [2]float64) float64 {
	return math.Sqrt((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]))
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
	{-2.95, -1.2},
	{-3.5, 0.7},
	{-0.8, 1.3},
	{1.7, 1.1},
	{-1.9,-1.9},
	/*{2.3,2},
	{2.6,1.5},
	{-1.9,2.3},
	{-2,1.3},
	{-4,1},
	{-3.5,-2.7},*/
}

func cast(s interface{}) [2]float64 {
	return s.([2]float64)
}
