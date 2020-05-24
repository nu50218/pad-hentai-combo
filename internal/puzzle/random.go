package puzzle

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomBoard(H, W, numColor int) Board {
	b := Board{
		Height: H,
		Width:  W,
		Data:   map[Coordinate]int{},
	}

	for x := 0; x < H; x++ {
		for y := 0; y < W; y++ {
			b.Data[Coordinate{X: x, Y: y}] = rand.Intn(numColor) + 1
		}
	}

	return b
}
