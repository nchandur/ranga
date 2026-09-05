package search

import (
	"context"
	"ranga/internal/board"
	"ranga/internal/evaluate/nnue"
)

// coordinates search tree execution
type Searcher struct {
	PV            PVTable                   // stores and tracks the pv line found during search
	TT            *TranspositionTable       // caches position evaluations and cutoffs
	Killers       [2][MAX_PLY]board.Move    // holds killer moves
	History       [12][64]int               // maintains history heuristic scores [piece][targetSq]
	Nodes         int                       // nodes visited that search
	NodeLimit     int                       // max number of nodes to visit
	*nnue.Network                           // reference to global weights
	Accumulators  [MAX_PLY]nnue.Accumulator // nnue accumulators for each ply
	Cancel        context.CancelFunc        // search cancel
}

// instantiates new searcher
func NewSearcher(nn *nnue.Network, ttSize int) *Searcher {
	s := Searcher{
		PV:        PVTable{},
		TT:        NewTranspositionTable(ttSize),
		Killers:   [2][MAX_PLY]board.Move{},
		History:   [12][64]int{},
		Network:   nn,
		Nodes:     0,
		NodeLimit: 0,
		Cancel:    nil,
	}
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

	staticEval := s.Network.Evaluate(&s.Accumulators[b.Ply], b.Side)

	// null move pruning (pass turn to attempt early fail-high)
	if score, prune := s.nullMovePruning(ctx, b, beta, depth, staticEval, inCheck); prune {
		return score
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
	movesSearched := 0

	localArray := [256]board.Move{}
	searchedQuiets := localArray[:0]

	for _, move := range ml.Moves[:ml.Count] {

		state := b.Preserve()
		s.UpdateAccumulator(b, move)
		b.Ply++

		if !b.MakeMove(move, false) {
			b.Ply--
			b.Restore(&state)
			continue
		}

		b.Repetition.Idx++
		b.Repetition.Table[b.Repetition.Idx] = b.Key

		legalMoves++

		// late move reduction
		score = s.lateMoveReduction(ctx, b, move, alpha, beta, depth, movesSearched, inCheck)

		b.Ply--
		b.Repetition.Idx--
		b.Restore(&state)
		movesSearched++

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

			// update pv line
			s.PV.updatePVLine(move, b.Ply)
		}

		if !move.IsCapture() {
			searchedQuiets = append(searchedQuiets, move)
		}

	}

	// checkmate or stalemate
	if legalMoves == 0 {
		var score int
		if inCheck {
			score = -ISMATE + b.Ply
		} else {
			score = 0
		}
		s.TT.Store(score, depth, b.Ply, FEXACT, b.Key, board.NOMOVE)
		return score
	}

	// store final score in transposition table
	s.TT.Store(alpha, depth, b.Ply, flag, b.Key, bestMove)
	return alpha
}

// executes search on a given state, returns the best move found
func (s *Searcher) Search(ctx context.Context, b *board.Board, depth int) (board.Move, int) {
	s.Reset()
	s.Accumulators[0].Refresh(s.Network, b)
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
		s.UpdateAccumulator(b, move)
		b.Ply++
		if !b.MakeMove(move, false) {
			b.Ply--
			b.Restore(&state)
			continue
		}

		// legal fallback in case of timeout
		if bestMove == board.NOMOVE {
			bestMove = move
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

	return bestMove, bestScore
}

// helper function to perform late move reduction
func (s *Searcher) lateMoveReduction(ctx context.Context, b *board.Board, move board.Move, alpha, beta, depth, movesSearched int, inCheck bool) int {

	var score int

	if movesSearched == 0 {
		score = -s.AlphaBeta(ctx, b, -beta, -alpha, depth-1)
	} else {
		reduced := movesSearched >= FULL_DEPTH_MOVES && depth >= REDUCTION_LIMIT &&
			!inCheck && !move.IsCapture() && move.Promoted() == board.Empty

		if reduced {
			score = -s.AlphaBeta(ctx, b, -alpha-1, -alpha, depth-2)
		} else {
			score = -s.AlphaBeta(ctx, b, -alpha-1, -alpha, depth-1)
		}

		// reduced search beat alpha
		if reduced && ctx.Err() == nil && score > alpha {
			score = -s.AlphaBeta(ctx, b, -alpha-1, -alpha, depth-1)
		}

		if ctx.Err() == nil && score > alpha && score < beta {
			score = -s.AlphaBeta(ctx, b, -beta, -alpha, depth-1)
		}
	}

	return score
}

// helper function for null move pruning
func (s *Searcher) nullMovePruning(ctx context.Context, b *board.Board, beta, depth, staticeval int, inCheck bool) (int, bool) {
	if depth < 3 ||
		inCheck ||
		b.Ply == 0 ||
		beta >= MATESCORE-MAX_PLY ||
		staticeval < beta ||
		!hasNonPawnMaterial(b) { // zugzwang check
		return 0, false
	}

	copy := b.Preserve()

	b.Side ^= 1

	b.Key ^= board.SideKey

	if b.EnPassant != board.NoSquare {
		b.Key ^= board.EnpassantKeys[b.EnPassant]
	}

	b.EnPassant = board.NoSquare

	b.Ply++

	// adaptive depth reduction
	R := 3 + depth/6
	if staticeval-beta > 200 {
		R++
	}
	reducedDepth := max(depth-R, 0)

	nullScore := -s.AlphaBeta(ctx, b, -beta, -beta+1, reducedDepth)
	b.Ply--

	b.Restore(&copy)

	if ctx.Err() != nil {
		return 0, false
	}

	// fail-high cutoff from null move
	if nullScore >= beta {
		if nullScore >= MATESCORE-MAX_PLY { // mate-range
			nullScore = beta
		}
		return nullScore, true
	}
	return 0, false
}
