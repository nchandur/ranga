package search

import "ranga/internal/board"

const (
	ScoreTT         int = 2_000_000
	ScorePV         int = 1_900_000
	ScoreQueenPromo int = 1_000_000
	ScoreCapture    int = 800_000
	ScoreKiller1    int = 600_000
	ScoreKiller2    int = 590_000
	ScoreUnderPromo int = -1_000_000
)

// evaluates heuristic score to move to prioritize it during move ordering
func (s *Searcher) scoreMove(b *board.Board, move, ttMove board.Move) int {

	if s.PV.ScorePV && move == s.PV.Table[0][b.Ply] {
		s.PV.ScorePV = false
		if move == ttMove {
			return ScoreTT
		}
		return ScorePV
	}

	if move == ttMove {
		return ScoreTT
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

		score := MVVLVA[b.Mailbox[move.Source()]][captured] + ScoreCapture

		// capturing into queen promotion is almost always the best capture
		if move.IsPromotion() && isQueenPromotion(b, move) {
			score += ScoreQueenPromo
		}

		return score

	} else {

		if move.IsPromotion() {
			if isQueenPromotion(b, move) {
				return ScoreQueenPromo
			}
			return ScoreUnderPromo
		}
		if s.stack[b.Ply].killers[0] == move {
			return ScoreKiller1
		} else if s.stack[b.Ply].killers[1] == move {
			return ScoreKiller2
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
