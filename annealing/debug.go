package annealing

import (
	"fmt"
	"io"
	"time"
)

func (a *Annealer) PrintInfo(w io.Writer) {
	a.mut.RLock()
	defer a.mut.RUnlock()

	fmt.Fprintln(w, "[INFO]", time.Now().Format(time.RFC3339), "cnt:", a.cnt, "best:", a.bestE)
	a.bestState.Debug()
}
