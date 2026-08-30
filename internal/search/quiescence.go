package search

import (
	"context"
	"ranga/internal/board"
)

// performs a quiescence search to evaluate quiet positions and avoid horizon effect
func (s *Searcher) Quiescence(ctx context.Context, b *board.Board, alpha, beta int) int {

	// guard against out-of-bounds at maximum search ply
	if b.Ply > MAX_PLY-1 {
		s.CheckAccumulatorDrift(b)
		return s.Network.Evaluate(&s.Accumulators[b.Ply], b.Side)
	}

	// check timeout or cancel
	if s.Nodes&2047 == 0 {

		if s.NodeLimit > 0 && s.Nodes >= s.NodeLimit {
			s.Cancel()
			return 0
		}

		if ctx.Err() != nil {
			return 0
		}
	}

	s.Nodes++

	// evaluate repetition or 50 move rule
	if (s.IsRepetition(b) || b.FiftyMove >= 100) && b.Ply != 0 {
		return 0
	}

	// check if current side king is in check
	var inCheck bool
	switch b.Side {
	case board.White:
		inCheck = b.IsSquareAttacked(board.Square(b.PieceBitBoards[board.WK].GetLSB()), board.Black)
	case board.Black:
		inCheck = b.IsSquareAttacked(board.Square(b.PieceBitBoards[board.BK].GetLSB()), board.White)
	}

	var standPat int

	// evaluate quiet position if not in check
	if !inCheck {
		standPat = s.Eval.Evaluate(b)

		// fail high
		if standPat >= beta {
			return beta
		}

		if standPat > alpha {
			alpha = standPat
		}
	}

	ml := board.NewMoveList()

	// generate escape moves when in check, otherwise search captures only
	if inCheck {
		ml.GenerateMoves(b)
	} else {
		ml.GenerateCaptures(b)
	}

	s.sortMove(b, ml, board.NOMOVE)

	legalMoves := 0

	for count := range ml.Count {

		copy := b.Preserve()
		s.UpdateAccumulator(b, ml.Moves[count])
		b.Ply++

		if !b.MakeMove(ml.Moves[count], false) {
			b.Ply--
			b.Restore(&copy)
			continue
		}

		b.Repetition.Idx++
		b.Repetition.Table[b.Repetition.Idx] = b.Key

		legalMoves++

		score := -s.Quiescence(ctx, b, -beta, -alpha)

		b.Ply--
		b.Repetition.Idx--

		b.Restore(&copy)

		if score > alpha {
			alpha = score

			if score >= beta {
				return beta
			}
		}

		// abort on context cancellation
		if ctx.Err() != nil {
			return 0
		}

	}

	// checkmate detection
	if inCheck && legalMoves == 0 {
		return -ISMATE + b.Ply
	}

	return alpha
}
