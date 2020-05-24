package puzzle

import "fmt"

type Board struct {
	Height int
	Width  int
	// 0 for empty
	Data map[Coordinate]int
}

func (b *Board) Simulate() Combo {
	initialCoordinate := map[Coordinate]Coordinate{}
	for x := 0; x < b.Height; x++ {
		for y := 0; y < b.Width; y++ {
			co := Coordinate{x, y}
			initialCoordinate[co] = co
		}
	}
	return b.simulate(0, initialCoordinate)
}

func (b *Board) simulate(phase int, initialCoordinate map[Coordinate]Coordinate) Combo {
	disappear := map[Coordinate]bool{}

	checkDisappear := func(c1, c2, c3 Coordinate) {
		if b.Data[c1] == 0 || b.Data[c2] == 0 || b.Data[c3] == 0 {
			return
		}

		if b.Data[c1] == b.Data[c2] && b.Data[c2] == b.Data[c3] {
			disappear[c1] = true
			disappear[c2] = true
			disappear[c3] = true
		}
	}

	for x := 0; x < b.Height-2; x++ {
		for y := 0; y < b.Width; y++ {
			c1 := Coordinate{x, y}
			c2 := Coordinate{x + 1, y}
			c3 := Coordinate{x + 2, y}

			checkDisappear(c1, c2, c3)
		}
	}

	for x := 0; x < b.Height; x++ {
		for y := 0; y < b.Width-2; y++ {
			c1 := Coordinate{x, y}
			c2 := Coordinate{x, y + 1}
			c3 := Coordinate{x, y + 2}

			checkDisappear(c1, c2, c3)
		}
	}

	if len(disappear) == 0 {
		// コンボ無し
		return nil
	}

	combo := Combo{}

	checked := map[Coordinate]bool{}
	for co := range disappear {
		if checked[co] {
			continue
		}

		chain := Chain{Phase: phase}

		// bfs
		queue := []Coordinate{co}
		checked[co] = true
		chain.Coordinates = append(chain.Coordinates, ComboCoordinate{
			co,
			initialCoordinate[co],
		})

		for ; len(queue) != 0; queue = queue[1:] {
			f := queue[0]

			dx := []int{0, 0, 1, -1}
			dy := []int{1, -1, 0, 0}

			for i := range dx {
				nextCo := Coordinate{f.X + dx[i], f.Y + dy[i]}
				if checked[nextCo] || b.Data[nextCo] == 0 {
					continue
				}
				if b.Data[co] != b.Data[nextCo] {
					continue
				}
				chain.Coordinates = append(chain.Coordinates, ComboCoordinate{
					nextCo,
					initialCoordinate[nextCo],
				})
				queue = append(queue, nextCo)
				checked[nextCo] = true
			}
		}

		combo = append(combo, chain)
	}

	// 残ったドロップを下に落とす
	nextBoard := Board{
		Height: b.Height,
		Width:  b.Width,
		Data:   map[Coordinate]int{},
	}
	nextInitialCoordinate := map[Coordinate]Coordinate{}
	for y := 0; y < b.Width; y++ {
		newX := b.Height - 1

		for x := b.Height - 1; 0 <= x; x-- {
			co := Coordinate{x, y}
			newCo := Coordinate{newX, y}

			if !disappear[co] {
				nextBoard.Data[newCo] = b.Data[co]
				nextInitialCoordinate[newCo] = initialCoordinate[co]
				newX--
			}
		}
	}

	return append(combo, nextBoard.simulate(phase+1, nextInitialCoordinate)...)
}

func (b *Board) Copy() Board {
	newBoard := Board{
		Height: b.Height,
		Width:  b.Width,
		Data:   make(map[Coordinate]int, len(b.Data)),
	}
	for key, value := range b.Data {
		newBoard.Data[key] = value
	}
	return newBoard
}

func (b *Board) Debug() {
	fmt.Println("[DEBUG] size :", b.Height, "*", b.Width)
	fmt.Println("[DEBUG] combo:", len(b.Simulate()))
	for x := 0; x < b.Height; x++ {
		fmt.Printf("[DEBUG] %02d| ", x+1)
		for y := 0; y < b.Width; y++ {
			fmt.Print(b.Data[Coordinate{x, y}])
			if y != b.Width-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
