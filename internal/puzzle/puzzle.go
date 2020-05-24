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
	label := map[Coordinate]int{}
	currentLabel := 1

	affixLabel := func(c1, c2, c3 Coordinate) {
		if b.Data[c1] == 0 || b.Data[c2] == 0 || b.Data[c3] == 0 {
			return
		}

		if b.Data[c1] == b.Data[c2] && b.Data[c2] == b.Data[c3] {
			// 既にChainの一部ならそのlabelを他に伝搬させる
			// 異なるlabelが存在することはない
			if label[c1] != 0 {
				label[c2] = label[c1]
				label[c3] = label[c1]
				return
			}
			if label[c2] != 0 {
				label[c1] = label[c2]
				label[c3] = label[c2]
				return
			}
			if label[c3] != 0 {
				label[c1] = label[c3]
				label[c2] = label[c3]
				return
			}

			// 新Chain
			label[c1] = currentLabel
			label[c2] = currentLabel
			label[c3] = currentLabel
			currentLabel++
		}
	}

	for x := 0; x < b.Height-2; x++ {
		for y := 0; y < b.Width; y++ {
			c1 := Coordinate{x, y}
			c2 := Coordinate{x + 1, y}
			c3 := Coordinate{x + 2, y}

			affixLabel(c1, c2, c3)
		}
	}

	for x := 0; x < b.Height; x++ {
		for y := 0; y < b.Width-2; y++ {
			c1 := Coordinate{x, y}
			c2 := Coordinate{x, y + 1}
			c3 := Coordinate{x, y + 2}

			affixLabel(c1, c2, c3)
		}
	}

	if currentLabel == 1 {
		// コンボ無し
		return nil
	}

	combo := Combo{}

	for l := 1; l < currentLabel; l++ {
		chain := Chain{Phase: phase}

		for x := 0; x < b.Height; x++ {
			for y := 0; y < b.Width; y++ {
				co := ComboCoordinate{
					Coordinate{x, y},
					initialCoordinate[Coordinate{x, y}],
				}
				if label[co.Coordinate] == l {
					chain.Coordinates = append(chain.Coordinates, co)
				}
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

			if label[co] == 0 {
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
