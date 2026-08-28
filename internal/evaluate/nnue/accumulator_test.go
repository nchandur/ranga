package nnue_test

import (
	"math/rand"
	"testing"

	"ranga/internal/board"
	"ranga/internal/nnue"
)

func startPosition() board.Board {
	b := board.NewBoard()
	b.ParseFEN(board.START)
	nnue.RefreshAll(&b)
	return b
}

func TestAccumulatorMatchesRefresh(t *testing.T) {
	nnue.Net = nnue.RandomNetwork()

	b := startPosition()
	rng := rand.New(rand.NewSource(1))

	for ply := range 60 {
		list := board.NewMoveList()
		list.GenerateMoves(&b)

		if list.Count == 0 {
			break
		}

		moved := false
		for _, i := range rng.Perm(list.Count) {
			if b.MakeMove(list.Moves[i], false) {
				moved = true
				break
			}
		}
		if !moved {
			break
		}

		checkAccumulator(t, &b, ply)
	}
}

func checkAccumulator(t *testing.T, b *board.Board, ply int) {
	t.Helper()

	incremental := b.NNUEAcc

	fresh := b.Preserve()
	nnue.RefreshAll(&fresh)

	const epsilon = 1e-3
	for p := range 2 {
		for i := range nnue.HiddenSize {
			diff := incremental[p][i] - fresh.NNUEAcc[p][i]
			if diff < -epsilon || diff > epsilon {
				t.Fatalf("ply %d: perspective %d feature %d mismatch: incremental=%f refresh=%f",
					ply, p, i, incremental[p][i], fresh.NNUEAcc[p][i])
			}
		}
	}
}
