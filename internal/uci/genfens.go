package uci

import (
	"math/rand/v2"
	"ranga/internal/board"
)

const (
	genfensPlies   = 8
	genfensEvalCap = 400 // reject openings scoring beyond +/- this many centipawns
	genfensDepth   = 4
)

// plays `plies` random legal moves from `start`, `count` times,
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
