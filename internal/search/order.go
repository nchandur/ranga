package search

import "ranga/internal/board"

// evaluates heuristic score to move to prioritize it during move ordering
func (s *Searcher) scoreMove(b *board.Board, move, ttMove board.Move) int {

	if move == ttMove {
		return 30000
	}

	if s.PV.ScorePV && move == s.PV.Table[0][b.Ply] {
		s.PV.ScorePV = false
		return 20000
	}

	if move.IsCapture() {
		var captured board.Piece

		if move.IsEnpass() {
			captured = board.WP
			if b.Side == board.White {
				captured = board.BP
			}
		} else {
			captured = b.Mailbox[move.Target()]
		}

		return MVVLVA[move.Piece()][captured] + 10000
	} else {

		if move.Promoted() == board.WQ || move.Promoted() == board.BQ {
			return 9500
		}

		if s.Killers[0][b.Ply] == move {
			return 9000
		} else if s.Killers[1][b.Ply] == move {
			return 8000
		} else {
			return s.History[move.Piece()][move.Target()]
		}

	}

}

// orders the legal moves in given move list
func (s *Searcher) sortMove(b *board.Board, ml *board.MoveList, ttMove board.Move) {

	score := make([]int, ml.Count)
	for i := 0; i < ml.Count; i++ {
		score[i] = s.scoreMove(b, ml.Moves[i], ttMove)
	}

	for i := 1; i < ml.Count; i++ {
		keyMove := ml.Moves[i]
		keyScore := score[i]
		j := i - 1

		for j >= 0 && score[j] < keyScore {
			ml.Moves[j+1] = ml.Moves[j]
			score[j+1] = score[j]
			j--
		}
		ml.Moves[j+1] = keyMove
		score[j+1] = keyScore
	}
}
