package search

import "ranga/internal/board"

// move ordering buckets
const (
	SCORE_TT        int = 2000000
	SCORE_PV        int = 1900000
	SCORE_CAPTURE   int = 1000000 // + MVVLVA
	SCORE_PROMOTION int = 950000
	SCORE_KILLER1   int = 900000
	SCORE_KILLER2   int = 800000
)

// evaluates heuristic score to move to prioritize it during move ordering
func (s *Searcher) scoreMove(b *board.Board, move, ttMove board.Move) int {

	if move == ttMove {
		return SCORE_TT
	}

	if s.PV.ScorePV && move == s.PV.Table[0][b.Ply] {
		s.PV.ScorePV = false
		return SCORE_PV
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

		return MVVLVA[b.Mailbox[move.Source()]][captured] + SCORE_CAPTURE
	} else {

		if move.IsPromotion() {
			return SCORE_PROMOTION
		}

		if s.stack[b.Ply].killers[0] == move {
			return SCORE_KILLER1
		} else if s.stack[b.Ply].killers[1] == move {
			return SCORE_KILLER2
		} else {
			return s.History[b.Mailbox[move.Source()]][move.Target()]
		}

	}

}

// orders the legal moves in given move list
func (s *Searcher) sortMove(b *board.Board, ml *board.MoveList, ttMove board.Move) {
	var scoreBuf [256]int
	score := scoreBuf[:ml.Count]

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
