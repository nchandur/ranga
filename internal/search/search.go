package search

import (
	"context"
	"ranga/internal/board"
	"ranga/internal/evaluate"
	"ranga/internal/evaluate/nnue"
)

// coordinates search tree execution
type Searcher struct {
	evaluate.Evaluator                        // static evaluation to score positions at leaf nodes
	PV                 PVTable                // stores and tracks the pv line found during search
	TT                 *TranspositionTable    // caches position evaluations and cutoffs
	Killers            [2][MAX_PLY]board.Move // holds killer moves
	History            [12][64]int            // maintains history heuristic scores [piece][targetSq]
	Nodes              int                    // nodes visited
	NodeLimit          int                    // cap for searched nodes
	Accumulators       [MAX_PLY]nnue.Accumulator
	Cancel             context.CancelFunc // cancel search
}

// instantiates new searcher
func NewSearcher(eval evaluate.Evaluator, ttSize int) *Searcher {
	s := Searcher{Evaluator: eval, PV: PVTable{}, TT: NewTranspositionTable(ttSize)}
	return &s
}

// clears searcher state
func (s *Searcher) Reset() {
	s.PV.Clear()
	s.Killers = [2][MAX_PLY]board.Move{}
	s.History = [12][64]int{}
}

// checks whether current board position has occurred previously
func (s *Searcher) IsRepetition(b *board.Board) bool {
	startIdx := max(b.Repetition.Idx-b.FiftyMove, 0)

	for i := startIdx; i < b.Repetition.Idx; i++ {
		if b.Repetition.Table[i] == b.Key {
			return true
		}
	}

	return false
}

// executes main alpha-beta minimax search tree traversal
func (s *Searcher) AlphaBeta(ctx context.Context, b *board.Board, alpha, beta, depth int) int {

	// guard against out-of-bounds at maximum search ply
	if b.Ply > MAX_PLY-1 {
		return s.Evaluate(b)
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

	s.PV.Length[b.Ply] = 0

	s.Nodes++

	// evaluate repetition or 50 move rule
	if (s.IsRepetition(b) || b.FiftyMove >= 100) && b.Ply != 0 {
		return 0
	}

	// transposition table lookup
	if score := s.TT.Probe(alpha, beta, b.Ply, depth, b.Key); b.Ply != 0 && score != NOENTRY && !s.PV.FollowPv {
		return score
	}

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
		return s.Quiescence(ctx, b, alpha, beta)
	}

	legalMoves := 0
	ml := board.NewMoveList()
	ml.GenerateMoves(b)

	// look for move in transposition table
	ttMove := s.TT.ProbeMove(b.Key)

	if s.PV.FollowPv {
		s.PV.enablePVScoring(ml, b.Ply)
	}

	s.sortMove(b, ml, ttMove)

	flag := FALPHA
	score := 0
	bestMove := board.NOMOVE

	localArray := [256]board.Move{}
	searchedQuiets := localArray[:0]

	for _, move := range ml.Moves[:ml.Count] {

		state := b.Preserve()
		b.Ply++

		if !b.MakeMove(move, false) {
			b.Ply--
			b.Restore(&state)
			continue
		}

		b.Repetition.Idx++
		b.Repetition.Table[b.Repetition.Idx] = b.Key

		legalMoves++

		score = -s.AlphaBeta(ctx, b, -beta, -alpha, depth-1)

		b.Ply--
		b.Repetition.Idx--
		b.Restore(&state)

		// abort on cancellation
		if ctx.Err() != nil {
			return 0
		}

		// fail hard on beta cutoff
		if score >= beta {
			s.TT.Store(score, depth, b.Ply, FBETA, b.Key, move)

			// record killer moves
			if !move.IsCapture() {
				s.Killers[1][b.Ply] = s.Killers[0][b.Ply]
				s.Killers[0][b.Ply] = move

				// record history heuristic
				bonus := min(depth*depth, 1200)
				s.updateHistory(move, bonus)

				// maluses for quiet moves already searched this node that didn't cut off
				for _, m := range searchedQuiets {
					s.updateHistory(m, -bonus)
				}

			}

			return beta
		}
		if score > alpha {
			alpha = score
			flag = FEXACT
			bestMove = move

			// history heuristic for quiet moves
			if !move.IsCapture() {
				bonus := min(depth*depth, 1200)
				s.updateHistory(move, bonus)
			}

			// update pv line
			s.PV.updatePVLine(move, b.Ply)
		}

		if !move.IsCapture() {
			searchedQuiets = append(searchedQuiets, move)
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

	// store final score in transposition table
	s.TT.Store(alpha, depth, b.Ply, flag, b.Key, bestMove)
	return alpha
}

// executes search on a given state, returns the best move found
func (s *Searcher) Search(ctx context.Context, b *board.Board, depth int) (board.Move, int, int) {

	s.Reset()
	alpha, beta := -INFINITY, INFINITY

	var bestMove board.Move
	bestScore := -INFINITY

	ml := board.NewMoveList()
	ml.GenerateMoves(b)

	s.PV.FollowPv = true
	s.PV.enablePVScoring(ml, 0)
	s.sortMove(b, ml, s.TT.ProbeMove(b.Key))

	for count, move := range ml.Moves[:ml.Count] {
		state := b.Preserve()
		b.Ply++
		if !b.MakeMove(move, false) {
			b.Ply--
			b.Restore(&state)
			continue
		}

		b.Repetition.Idx++
		b.Repetition.Table[b.Repetition.Idx] = b.Key

		s.PV.FollowPv = (count == 0)
		score := -s.AlphaBeta(ctx, b, -beta, -alpha, depth-1)

		b.Ply--
		b.Repetition.Idx--
		b.Restore(&state)

		if ctx.Err() != nil {
			break
		}

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if score > alpha {
			alpha = score
			s.PV.updatePVLine(move, 0)
		}

	}

	if s.PV.Length[0] > 0 {
		bestMove = s.PV.Table[0][0]
	}

	return bestMove, bestScore, s.Nodes
}
