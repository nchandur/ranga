package search

import "ranga/internal/board"

// update history heuristic with gravity
func (s *Searcher) updateHistory(move board.Move, bonus int) {
	clamped := clamp(bonus, -MAX_HISTORY, MAX_HISTORY)
	s.History[move.Piece()][move.Target()] += (clamped - s.History[move.Piece()][move.Target()]*abs(clamped)/MAX_HISTORY)
}

// returns true for side with non-pawn material
func hasNonPawnMaterial(b *board.Board) bool {

	switch b.Side {
	case board.White:
		return b.PieceBitBoards[board.WN].CountBits() > 0 ||
			b.PieceBitBoards[board.WB].CountBits() > 0 ||
			b.PieceBitBoards[board.WR].CountBits() > 0 ||
			b.PieceBitBoards[board.WQ].CountBits() > 0
	case board.Black:
		return b.PieceBitBoards[board.BN].CountBits() > 0 ||
			b.PieceBitBoards[board.BB].CountBits() > 0 ||
			b.PieceBitBoards[board.BR].CountBits() > 0 ||
			b.PieceBitBoards[board.BQ].CountBits() > 0
	}

	return false
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
