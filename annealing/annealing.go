package annealing

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/nu50218/pad-hentai-combo/internal/puzzle"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type EvalFunc func(puzzle.Board) int
type NeighbourFunc func(puzzle.Board) puzzle.Board

type Annealer struct {
	eval      EvalFunc
	neighbour NeighbourFunc
	bestState puzzle.Board
	bestE     int
	mut       sync.RWMutex
	option    Option
	cnt       int
}

type Option struct {
	H        int
	W        int
	NumColor int
	// alpha 温度の底
	Alpha float64
	// 盤面最大コンボの最小数
	MinMaxCombo int
}

var DefaultOption = Option{
	H:           5,
	W:           6,
	NumColor:    6,
	Alpha:       0.995,
	MinMaxCombo: 9,
}

func NewAnnealer(eval EvalFunc, neighbour NeighbourFunc, option Option) *Annealer {
	tmpState := puzzle.RandomBoard(option.H, option.W, option.NumColor)

	return &Annealer{
		eval:      eval,
		neighbour: neighbour,
		bestState: tmpState,
		bestE:     eval(tmpState),
		option:    option,
	}
}

func (a *Annealer) SimulatedAnnealing(maxIter int) {
	currentState := puzzle.RandomBoard(a.option.H, a.option.W, a.option.NumColor)

	countMaxCombo := func(b puzzle.Board) int {
		// color -> num
		cnt := map[int]int{}
		for x := 0; x < b.Height; x++ {
			for y := 0; y < b.Width; y++ {
				cnt[b.Data[puzzle.Coordinate{X: x, Y: y}]]++
			}
		}

		maxCombo := 0
		for _, val := range cnt {
			maxCombo += val / 3
		}

		return maxCombo
	}

	for countMaxCombo(currentState) < a.option.MinMaxCombo {
		currentState = puzzle.RandomBoard(a.option.H, a.option.W, a.option.NumColor)
	}

	a.mut.Lock()
	e := a.eval(currentState)
	if a.bestE < e {
		a.bestState = currentState.Copy()
		a.bestE = e
	}
	a.mut.Unlock()

	for iter := 0; iter < maxIter; iter++ {
		nextState := a.neighbour(currentState)
		nextE := a.eval(nextState)

		a.mut.Lock()
		if a.bestE < nextE {
			a.bestState = nextState.Copy()
			a.bestE = nextE
		}
		a.mut.Unlock()

		temp := a.temperature(float64(iter) / float64(maxIter))
		if rand.Float64() <= a.probability(e, nextE, temp) {
			currentState = nextState
			e = nextE
		}
	}

	a.cnt++
}

func (*Annealer) probability(e1, e2 int, t float64) float64 {
	if e1 <= e2 {
		return 1.0
	}
	return math.Pow(math.E, float64(e2-e1)/t)
}

func (a *Annealer) temperature(r float64) float64 {
	return math.Pow(a.option.Alpha, r)
}
