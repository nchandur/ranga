package search

import (
	"context"
	"ranga/internal/board"
	"testing"
)

func TestQuiescence(t *testing.T) {
	eval := &mockEvaluator{}
	tt := NewTranspositionTable(DEFAULT_TT_SIZE)

	t.Run("stand-pat fails high when eval >= beta", func(t *testing.T) {
		s := NewSearcher(eval, tt)
		b := board.NewBoard()
		b.ParseFEN(board.START)

		alpha := -100
		beta := -50

		score := s.Quiescence(context.Background(), &b, alpha, beta)
		if score != beta {
			t.Errorf("expected stand-pat beta cutoff %d, got %d", beta, score)
		}
	})
	t.Run("stand-pat updates alpha when eval > alpha", func(t *testing.T) {
		s := NewSearcher(eval, tt)
		b := board.NewBoard()
		b.ParseFEN("4k3/8/8/8/8/8/8/Q3K3 w - - 0 1")

		alpha := 0
		beta := 2000

		score := s.Quiescence(context.Background(), &b, alpha, beta)
		if score < 900 {
			t.Errorf("expected score to reflect at least the Queen advantage, got %d", score)
		}
	})
	t.Run("resolves tactical captures instead of standing pat", func(t *testing.T) {
		s := NewSearcher(eval, tt)
		b := board.NewBoard()
		b.ParseFEN("4k3/8/8/3q4/4P3/8/8/4K3 w - - 0 1")

		alpha := -2000
		beta := 2000

		score := s.Quiescence(context.Background(), &b, alpha, beta)
		if score <= 0 {
			t.Errorf("Quiescence failed to resolve hanging queen capture: score %d", score)
		}
	})
	t.Run("searches quiet moves to escape check", func(t *testing.T) {
		s := NewSearcher(eval, tt)
		b := board.NewBoard()
		b.ParseFEN("4r3/8/8/8/8/8/8/4K2k w - - 0 1")

		alpha := -INFINITY
		beta := INFINITY

		score := s.Quiescence(context.Background(), &b, alpha, beta)

		if score <= -ISMATE+100 {
			t.Errorf("Quiescence erroneously evaluated check evasion as checkmate: score %d", score)
		}
	})
	t.Run("detects checkmate when in check with no legal evasions", func(t *testing.T) {
		s := NewSearcher(eval, tt)
		b := board.NewBoard()
		b.ParseFEN("k7/8/8/8/8/8/1r6/r3K3 w - - 0 1")
		b.Ply = 2

		score := s.Quiescence(context.Background(), &b, -INFINITY, INFINITY)

		wantScore := -ISMATE + b.Ply
		if score != wantScore {
			t.Errorf("expected checkmate score %d, got %d", wantScore, score)
		}
	})
	t.Run("deltaPruning identifies futile captures", func(t *testing.T) {
		s := NewSearcher(eval, tt)
		b := board.NewBoard()
		b.ParseFEN("4k3/8/8/3p4/4P3/8/8/4K3 w - - 0 1")

		move := board.NewMove(board.E4, board.D5, board.WP, board.Empty, true, false, false, false)

		evaluation := -2000
		alpha := 500
		if !s.deltaPruning(alpha, evaluation, false, &b, move) {
			t.Errorf("expected futile capture to be pruned by deltaPruning")
		}

		if s.deltaPruning(alpha, 400, false, &b, move) {
			t.Errorf("expected viable capture not to be pruned")
		}

		if s.deltaPruning(alpha, evaluation, true, &b, move) {
			t.Errorf("deltaPruning should not prune when in check")
		}
	})
}
