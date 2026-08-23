package search

import (
	"context"
	"ranga/internal/board"
)

// performs a quiescence search to evaluate quiet positions and avoid horizon effect
func (s *Searcher) Quiescence(ctx context.Context, b *board.Board, alpha, beta int, nodes *int) int {

	// guard against out-of-bounds at maximum search ply
	if b.Ply > MAX_PLY-1 {
		return s.Evaluate(b)
	}

	// check timeout or cancel
	if (*nodes)&2047 == 0 {
		if ctx.Err() != nil {
			return 0
		}
	}

	(*nodes)++

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
		standPat = s.Evaluate(b)

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

	s.sortMove(b, ml)

	legalMoves := 0

	for count := range ml.Count {

		copy := b.Preserve()
		b.Ply++

		if !b.MakeMove(ml.Moves[count], false) {
			b.Ply--
			b.Restore(&copy)
			continue
		}

		legalMoves++

		score := -s.Quiescence(ctx, b, -beta, -alpha, nodes)

		b.Ply--

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
