package search

import "ranga/internal/board"

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
