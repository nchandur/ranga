package uci

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"ranga/internal/board"
	"sync/atomic"
)

const (
	genfensPlies   = 8
	genfensEvalCap = 400  // reject openings scoring beyond +/- this many centipawns
	genfensNodes   = 5000 // node budget for scoring each candidate
)

// generates `count` random openings from the current position
func (e *Engine) generate(count int, seed uint64, emit func(*board.Board)) int {
	e.pauseSearch()
	stop := new(atomic.Bool)
	e.stop = stop
	e.searcher.Stop = stop

	accept := func(b *board.Board) bool {
		return abs(e.scorePosition(b)) <= genfensEvalCap
	}

	n := genPositions(&e.board, count, genfensPlies, seed, accept, emit)

	e.searcher.NN.Reset(&e.board)
	e.searcher.TT.Clear()

	return n
}

// generates `count` positions and saves just the FENs
func (e *Engine) GenFENsToFile(path string, count int, seed uint64) (int, error) {
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	w := bufio.NewWriter(f)

	n := e.generate(count, seed, func(b *board.Board) {
		fmt.Fprintln(w, b.FEN())
	})

	if err := w.Flush(); err != nil {
		f.Close()
		return n, err
	}
	if err := f.Close(); err != nil {
		return n, err
	}
	if n < count {
		return n, fmt.Errorf("produced %d of %d positions", n, count)
	}
	return n, nil
}

// plays `plies` random legal moves from `start`, `count` times
func genPositions(start *board.Board, count, plies int, seed uint64,
	accept func(*board.Board) bool, emit func(*board.Board)) int {

	rng := rand.New(rand.NewPCG(seed, seed))

	produced := 0
	for attempts := 0; produced < count && attempts < 1000*count; attempts++ {
		b := start.Preserve()

		ok := true
		for range plies {
			if !playRandomMove(&b, rng) {
				ok = false
				break
			}
		}
		if !ok || !hasLegalMove(&b) {
			continue
		}
		if accept != nil && !accept(&b) {
			continue
		}

		emit(&b)
		produced++
	}
	return produced
}

// plays one random legal move
func playRandomMove(b *board.Board, rng *rand.Rand) bool {
	list := board.NewMoveList()
	list.GenerateMoves(b)

	for _, i := range rng.Perm(list.Count) {
		if b.MakeMove(list.Moves[i], false) {
			b.Ply++
			return true
		}
	}
	return false
}

// reports whether side to move has at least one legal move
func hasLegalMove(b *board.Board) bool {
	list := board.NewMoveList()
	list.GenerateMoves(b)

	for i := range list.Count {
		c := b.Preserve()
		if c.MakeMove(list.Moves[i], false) {
			return true
		}
	}
	return false
}
