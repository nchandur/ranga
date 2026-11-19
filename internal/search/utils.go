package search

import (
	"ranga/internal/board"
)

// update main + continuation history with gravity
func (s *Searcher) updateHistory(b *board.Board, move board.Move, bonus int) {
	piece, to := b.Mailbox[move.Source()], move.Target()
	clamped := clamp(bonus, -MAX_HISTORY, MAX_HISTORY)

	h := &s.History[piece][to]
	*h += clamped - *h*abs(clamped)/MAX_HISTORY

	if ch := s.contEntry(b.Ply, 1); ch != nil {
		c := &ch[piece][to]
		*c += int16(clamped - int(*c)*abs(clamped)/MAX_HISTORY)
	}
}

// update killer move heuristic
func (s *Searcher) updateKiller(ply int, m board.Move) {
	n := &s.stack[ply]

	if n.killers[0] == m {
		return
	}

	n.killers[1] = n.killers[0]
	n.killers[0] = m
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

// reports whether move promotes to queen for side to move
func isQueenPromotion(b *board.Board, move board.Move) bool {
	p := move.Promoted(b.Side)
	return p == board.WQ || p == board.BQ
}

// returns the stored refutation of the previous move, or NOMOVE
func (s *Searcher) counterMove(ply int) board.Move {
	if ply == 0 {
		return board.NOMOVE
	}
	prev := &s.stack[ply-1]
	if prev.currMove == board.NOMOVE {
		return board.NOMOVE
	}
	return s.CounterMoves[prev.currPiece][prev.currMove.Target()]
}

func (s *Searcher) updateCounter(ply int, m board.Move) {
	if ply == 0 {
		return
	}
	prev := &s.stack[ply-1]
	if prev.currMove == board.NOMOVE {
		return
	}
	s.CounterMoves[prev.currPiece][prev.currMove.Target()] = m
}

// table slice keyed by the move played `offset` plies ago, nil if none (root or null move)
func (s *Searcher) contEntry(ply, offset int) *[12][64]int16 {
	if ply < offset {
		return nil
	}
	n := &s.stack[ply-offset]
	if n.currMove == board.NOMOVE {
		return nil
	}
	return &s.ContHist[n.currPiece][n.currMove.Target()]
}
