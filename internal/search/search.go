package search

import (
	"ranga/internal/board"
	"ranga/internal/evaluate"
	"ranga/internal/evaluate/nnue"
	"sync/atomic"
)

// coordinates search tree execution
type Searcher struct {
	evaluate.Evaluator // static evaluation to score positions at leaf nodes
	NN                 *nnue.NNUE
	PV                 PVTable                // stores and tracks the pv line found during search
	TT                 *TranspositionTable    // caches position evaluations and cutoffs
	Killers            [2][MAX_PLY]board.Move // holds killer moves
	History            [12][64]int            // maintains history heuristic scores [piece][targetSq]
	Nodes              int                    // nodes visited that search
	NodeLimit          int                    // max number of nodes to visit
	NodesPub           atomic.Uint64          // last published node count
	Stop               *atomic.Bool           // shared stop flag across all threads
}

// instantiates new searcher
func NewSearcher(eval evaluate.Evaluator, tt *TranspositionTable, stop *atomic.Bool) *Searcher {
	s := Searcher{
		Evaluator: eval,
		PV:        PVTable{},
		TT:        tt,
		Killers:   [2][MAX_PLY]board.Move{},
		History:   [12][64]int{},
		Nodes:     0,
		NodeLimit: 0,
		Stop:      stop,
	}

	// only evaluator is NNUE
	if nn, ok := eval.(*nnue.NNUE); ok {
		s.NN = nn
	}

	return &s
}

func (s *Searcher) NewHelper() *Searcher {
	var eval evaluate.Evaluator = s.Evaluator
	if s.NN != nil {
		eval = s.NN.Clone()
	}
	return NewSearcher(eval, s.TT, s.Stop)
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
func (s *Searcher) AlphaBeta(b *board.Board, alpha, beta, depth int) int {
	// guard against out-of-bounds at maximum search ply
	if b.Ply >= MAX_PLY-1 || b.Repetition.Idx >= len(b.Repetition.Table)-1 {
		return s.Evaluate(b)
	}

	// check timeout or cancel
	if s.Nodes&2047 == 0 {
		s.NodesPub.Store(uint64(s.Nodes))
		if s.NodeLimit > 0 && s.Nodes >= s.NodeLimit {
			s.Stop.Load()
			return 0
		}

		if s.stopped() {
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
		return s.Quiescence(b, alpha, beta)
	}

	staticEval := s.Evaluate(b)
	pvNode := beta-alpha > 1

	// reverse futility pruning
	if score, prune := s.reverseFutilityPruning(beta, depth, staticEval, inCheck, pvNode); prune {
		return score
	}

	// null move pruning (pass turn to attempt early fail-high)
	if score, prune := s.nullMovePruning(b, beta, depth, staticEval, inCheck); prune {
		return score
	}

	// futility pruning (discard moves with no potential of improving alpha)
	isFP := s.isFutile(alpha, depth, staticEval, inCheck, pvNode)

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

		// check move futility
		if isFP &&
			movesSearched > 0 && // search at least one move
			!move.IsCapture() && move.Promoted() == board.Empty {
			continue
		}

		state := b.Preserve()

		var nnState nnue.Snapshot
		if s.NN != nil {
			nnState = s.NN.Preserve()
		}

		b.Ply++

		if !b.MakeMove(move, false) {
			b.Ply--
			b.Restore(&state)
			continue
		}

		if s.NN != nil {
			s.NN.Update(&state, move)
		}

		b.Repetition.Idx++
		b.Repetition.Table[b.Repetition.Idx] = b.Key

		legalMoves++

		// late move reduction
		score = s.lateMoveReduction(b, move, alpha, beta, depth, movesSearched, inCheck)

		b.Ply--
		b.Repetition.Idx--
		b.Restore(&state)

		if s.NN != nil {
			s.NN.Restore(nnState)
		}

		movesSearched++

		// abort on cancellation
		if s.stopped() {
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
func (s *Searcher) Search(b *board.Board, depth int) (board.Move, int) {
	s.Reset()

	if s.NN != nil {
		s.NN.Reset(b)
	}

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

		var nnState nnue.Snapshot
		if s.NN != nil {
			nnState = s.NN.Preserve()
		}

		b.Ply++
		if !b.MakeMove(move, false) {
			b.Ply--
			b.Restore(&state)
			continue
		}

		if s.NN != nil {
			s.NN.Update(&state, move)
		}

		// legal fallback in case of timeout
		if bestMove == board.NOMOVE {
			bestMove = move
		}

		b.Repetition.Idx++
		b.Repetition.Table[b.Repetition.Idx] = b.Key

		s.PV.FollowPv = (count == 0)
		score := -s.AlphaBeta(b, -beta, -alpha, depth-1)

		b.Ply--
		b.Repetition.Idx--
		b.Restore(&state)

		if s.NN != nil {
			s.NN.Restore(nnState)
		}

		if s.stopped() {
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
func (s *Searcher) lateMoveReduction(b *board.Board, move board.Move, alpha, beta, depth, movesSearched int, inCheck bool) int {

	var score int

	if movesSearched == 0 {
		score = -s.AlphaBeta(b, -beta, -alpha, depth-1)
	} else {
		reduced := movesSearched >= FULL_DEPTH_MOVES && depth >= REDUCTION_LIMIT &&
			!inCheck && !move.IsCapture() && move.Promoted() == board.Empty

		if reduced {
			score = -s.AlphaBeta(b, -alpha-1, -alpha, depth-2)
		} else {
			score = -s.AlphaBeta(b, -alpha-1, -alpha, depth-1)
		}

		// reduced search beat alpha
		if reduced && !s.stopped() && score > alpha {
			score = -s.AlphaBeta(b, -alpha-1, -alpha, depth-1)
		}

		if !s.stopped() && score > alpha && score < beta {
			score = -s.AlphaBeta(b, -beta, -alpha, depth-1)
		}
	}

	return score
}

// helper function for null move pruning
func (s *Searcher) nullMovePruning(b *board.Board, beta, depth, staticeval int, inCheck bool) (int, bool) {
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

	nullScore := -s.AlphaBeta(b, -beta, -beta+1, reducedDepth)
	b.Ply--

	b.Restore(&copy)

	if s.stopped() {
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

// helper function for futlitity pruning
func (s *Searcher) isFutile(alpha, depth, staticEval int, inCheck, pvNode bool) bool {
	if depth > 3 || inCheck || pvNode || alpha <= -MATESCORE+MAX_PLY || alpha >= MATESCORE-MAX_PLY {
		return false
	}

	return staticEval+FutilityMargin[depth] <= alpha
}

// helper function for reverse futility pruning
func (s *Searcher) reverseFutilityPruning(beta, depth, staticEval int, inCheck, isPVNode bool) (int, bool) {

	if depth >= 8 || inCheck || isPVNode || beta <= -MATESCORE+MAX_PLY || beta >= MATESCORE-MAX_PLY {
		return 0, false
	}

	if staticEval-(150*depth) >= beta {
		return staticEval, true
	}

	return 0, false
}
