package search

import (
	"context"
	"ranga/internal/board"
	"ranga/internal/evaluate"
)

// coordinates search tree execution
type Searcher struct {
	evaluate.Evaluator         // static evaluation to score positions at leaf nodes
	PV                 PVTable // stores and tracks the pv line found during search
}

// instantiates new searcher
func NewSearcher(eval evaluate.Evaluator) *Searcher {
	s := Searcher{Evaluator: eval}
	return &s
}

// clears searcher state
func (s *Searcher) Reset() {
	s.PV.Clear()
}

// executes main alpha-beta minimax search tree traversal
func (s *Searcher) AlphaBeta(ctx context.Context, b *board.Board, alpha, beta, depth int, nodes *int) int {
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

	s.PV.Length[b.Ply] = 0

	(*nodes)++

	var inCheck bool

	// check if current side king is in check
	switch b.Side {
	case board.White:
		inCheck = b.IsSquareAttacked(board.Square(b.PieceBitBoards[board.WK].GetLSB()), board.Black)
	case board.Black:
		inCheck = b.IsSquareAttacked(board.Square(b.PieceBitBoards[board.BK].GetLSB()), board.White)
	}

	// increase depth by 1 if in check
	if inCheck {
		depth++
	}

	// drop into quiescence search at leaf nodes
	if depth == 0 {
		return s.Quiescence(ctx, b, alpha, beta, nodes)
	}

	legalMoves := 0
	ml := board.NewMoveList()
	ml.GenerateMoves(b)

	if s.PV.FollowPv {
		s.PV.enablePVScoring(ml, b.Ply)
	}

	s.sortMove(b, ml)

	for _, move := range ml.Moves[:ml.Count] {

		state := b.Preserve()
		b.Ply++
		if !b.MakeMove(move, false) {
			b.Restore(&state)
			b.Ply--
			continue
		}

		legalMoves++

		score := -s.AlphaBeta(ctx, b, -beta, -alpha, depth-1, nodes)

		b.Restore(&state)
		b.Ply--

		// abort on cancellation
		if ctx.Err() != nil {
			return 0
		}

		// fail hard on beta cutoff
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score

			// update pv line
			s.PV.updatePVLine(move, b.Ply)
		}

	}

	// checkmate or stalemate
	if legalMoves == 0 {
		if inCheck {
			return -ISMATE + b.Ply
		} else {
			return 0
		}
	}

	return alpha
}

// executes search on a given state, returns the best move found
func (s *Searcher) Search(ctx context.Context, b *board.Board, depth int) (board.Move, int, int) {
	nodes := 0
	s.PV.Clear()
	alpha, beta := -INFINITY, INFINITY

	var bestMove board.Move
	bestScore := -INFINITY

	ml := board.NewMoveList()
	ml.GenerateMoves(b)

	s.PV.FollowPv = true
	s.PV.enablePVScoring(ml, 0)
	s.sortMove(b, ml)

	for _, move := range ml.Moves[:ml.Count] {
		state := b.Preserve()
		b.Ply++
		if !b.MakeMove(move, false) {
			b.Restore(&state)
			b.Ply--
			continue
		}

		s.PV.FollowPv = true
		score := -s.AlphaBeta(ctx, b, -beta, -alpha, depth-1, &nodes)

		b.Restore(&state)
		b.Ply--

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if score > alpha {
			alpha = score
			s.PV.updatePVLine(move, 0)
		}

		if ctx.Err() != nil {
			break
		}

	}

	if s.PV.Length[0] > 0 {
		bestMove = s.PV.Table[0][0]
	}

	return bestMove, bestScore, nodes
}
