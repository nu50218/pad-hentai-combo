package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nu50218/pad-hentai-combo/annealing"
	"github.com/nu50218/pad-hentai-combo/internal/puzzle"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func eval(b puzzle.Board) int {
	combo := b.Simulate()

	const (
		k1 = 0
		k2 = 0
		k3 = 1
	)

	// コンボ加点
	score := k1 * len(combo)

	// もともと離れていればいるほど加点
	dist := func(co1, co2 puzzle.ComboCoordinate) int {
		diff := func(x, y int) int {
			if x < y {
				return y - x
			}
			return x - y
		}
		return diff(co1.Initial.X, co2.Initial.X) + diff(co1.Initial.Y, co2.Initial.Y)
	}

	for _, chain := range combo {
		maxDist := 0

		for _, co1 := range chain.Coordinates {
			for _, co2 := range chain.Coordinates {
				if maxDist < dist(co1, co2) {
					maxDist = dist(co1, co2)
				}
			}
		}

		score += k2 * (maxDist - 2)
	}

	// 落とし加点
	for _, chain := range combo {
		score += k3 * chain.Phase * chain.Phase
	}

	return score
}

func neighbour(b puzzle.Board) puzzle.Board {
	randomSwap := func() {
		co1 := puzzle.Coordinate{
			X: rand.Intn(b.Height),
			Y: rand.Intn(b.Width),
		}
		co2 := puzzle.Coordinate{
			X: rand.Intn(b.Height),
			Y: rand.Intn(b.Width),
		}

		b.Data[co1], b.Data[co2] = b.Data[co2], b.Data[co1]
	}

	randomSwap()
	randomSwap()
	randomSwap()

	return b
}

func main() {
	option := annealing.DefaultOption
	option.Alpha = 0.999995
	a := annealing.NewAnnealer(eval, neighbour, option)

	go func() {
		for {
			a.PrintInfo(os.Stderr)
			time.Sleep(1 * time.Second)
		}
	}()

	cap := make(chan struct{}, 4)
	for {
		cap <- struct{}{}
		go func() {
			a.SimulatedAnnealing(100000)
			<-cap
		}()
	}
}
