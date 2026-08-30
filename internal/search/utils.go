package search

import (
	"fmt"
	"ranga/internal/board"
	"ranga/internal/evaluate/nnue"
)

// update history heuristic with gravity
func (s *Searcher) updateHistory(move board.Move, bonus int) {
	clamped := clamp(bonus, -MAX_HISTORY, MAX_HISTORY)
	s.History[move.Piece()][move.Target()] += (clamped - s.History[move.Piece()][move.Target()]*abs(clamped)/MAX_HISTORY)
}

// helper function to clamp bonus between MAX_HISTORY
func clamp(n, low, high int) int {
	if n <= low {
		return low
	}

	if n >= high {
		return high
	}

	return n
}

// helper function to calculate absolute value of an integer
func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

// copies the current ply's accumulator to next ply and applies move delta
func (s *Searcher) UpdateAccumulator(b *board.Board, move board.Move) {
	nextPly := b.Ply + 1

	// copy current state to next ply
	s.Accumulators[nextPly] = s.Accumulators[b.Ply]
	acc := &s.Accumulators[nextPly]

	source := move.Source()
	target := move.Target()
	movingPiece := b.Mailbox[source]

	// remove moving piece from its original square
	acc.RemovePiece(s.Network, movingPiece, source)

	// remove captured piece from target
	if move.IsCapture() && !move.IsEnpass() {
		capturedPiece := b.Mailbox[target]
		acc.RemovePiece(s.Network, capturedPiece, target)
	}

	// handle enpassant capture
	if move.IsEnpass() {
		if b.Side == board.White {
			acc.RemovePiece(s.Network, board.BP, target+8)
		} else {
			acc.RemovePiece(s.Network, board.WP, target-8)
		}
	}

	// handle promotions
	if move.Promoted() != board.Empty {
		acc.AddPiece(s.Network, move.Promoted(), target)
	} else {
		acc.AddPiece(s.Network, movingPiece, target)
	}

	// handle castling
	if move.IsCastle() {
		switch target {
		case board.G1:
			acc.RemovePiece(s.Network, board.WR, board.H1)
			acc.AddPiece(s.Network, board.WR, board.F1)
		case board.C1:
			acc.RemovePiece(s.Network, board.WR, board.A1)
			acc.AddPiece(s.Network, board.WR, board.D1)
		case board.G8:
			acc.RemovePiece(s.Network, board.BR, board.H8)
			acc.AddPiece(s.Network, board.BR, board.F8)
		case board.C8:
			acc.RemovePiece(s.Network, board.BR, board.A8)
			acc.AddPiece(s.Network, board.BR, board.D8)
		}
	}
}

// verifies that the incrementally updated accumulator
func (s *Searcher) CheckAccumulatorDrift(b *board.Board) {
	// 1. Calculate a fresh accumulator from scratch
	var fresh nnue.Accumulator
	fresh.Refresh(s.Network, b)

	// 2. Get the current incrementally updated accumulator
	current := &s.Accumulators[b.Ply]

	// 3. Compare all 256 features for both White and Black
	for i := range nnue.HiddenSize {
		if fresh.White[i] != current.White[i] || fresh.Black[i] != current.Black[i] {
			fmt.Printf("\n--- ACCUMULATOR DRIFT DETECTED ---\n")
			fmt.Printf("Ply: %d\n", b.Ply)
			fmt.Printf("Mismatch at index %d\n", i)
			fmt.Printf("White - Fresh: %d | Incremental: %d\n", fresh.White[i], current.White[i])
			fmt.Printf("Black - Fresh: %d | Incremental: %d\n", fresh.Black[i], current.Black[i])

			b.Print()

			panic("Accumulator drift detected! Halting search.")
		}
	}
}
