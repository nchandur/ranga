package search

import (
	"ranga/internal/board"
	"testing"
)

func TestSearchUtils(t *testing.T) {
	t.Run("abs returns absolute value", func(t *testing.T) {
		tests := []struct {
			input int
			want  int
		}{
			{input: 0, want: 0},
			{input: 15, want: 15},
			{input: -25, want: 25},
			{input: -10000, want: 10000},
		}

		for _, tt := range tests {
			if got := abs(tt.input); got != tt.want {
				t.Errorf("abs(%d) = %d; want %d", tt.input, got, tt.want)
			}
		}
	})
	t.Run("clamp bounds values correctly", func(t *testing.T) {
		tests := []struct {
			n, low, high int
			want         int
		}{
			{n: 5, low: 0, high: 10, want: 5},
			{n: -5, low: 0, high: 10, want: 0},
			{n: 15, low: 0, high: 10, want: 10},
			{n: 0, low: 0, high: 10, want: 0},
			{n: 10, low: 0, high: 10, want: 10},
		}

		for _, tt := range tests {
			if got := clamp(tt.n, tt.low, tt.high); got != tt.want {
				t.Errorf("clamp(%d, %d, %d) = %d; want %d", tt.n, tt.low, tt.high, got, tt.want)
			}
		}
	})
	t.Run("hasNonPawnMaterial detects presence of major/minor pieces", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1")
		if hasNonPawnMaterial(&b) {
			t.Errorf("expected false for White in pawn endgame")
		}

		b.Side = board.Black
		if hasNonPawnMaterial(&b) {
			t.Errorf("expected false for Black in pawn endgame")
		}

		b.ParseFEN("4k3/pppppppp/8/8/8/8/PPPPPPPP/1N2K3 w - - 0 1")
		if !hasNonPawnMaterial(&b) {
			t.Errorf("expected true for White with Knight on b1")
		}

		b.Side = board.Black
		if hasNonPawnMaterial(&b) {
			t.Errorf("expected false for Black despite White possessing a Knight")
		}

		b.ParseFEN("3qk3/8/8/8/8/8/8/4K3 b - - 0 1")
		if !hasNonPawnMaterial(&b) {
			t.Errorf("expected true for Black with Queen on d8")
		}
	})
	t.Run("updateHistory updates score with gravity scaling", func(t *testing.T) {
		s := &Searcher{}
		move := board.NewMove(board.E2, board.E4, board.WP, board.Empty, false, false, false, false)

		piece := move.Piece()
		target := move.Target()

		if s.History[piece][target] != 0 {
			t.Fatalf("expected initial history 0, got %d", s.History[piece][target])
		}

		bonus := 100
		s.updateHistory(move, bonus)

		firstScore := s.History[piece][target]
		if firstScore <= 0 {
			t.Errorf("expected history score to increase, got %d", firstScore)
		}

		s.updateHistory(move, bonus)
		secondScore := s.History[piece][target]
		if secondScore <= firstScore {
			t.Errorf("expected history score to increase further, got %d (prev: %d)", secondScore, firstScore)
		}

		hugeBonus := MAX_HISTORY * 5
		s.updateHistory(move, hugeBonus)
		if s.History[piece][target] > MAX_HISTORY {
			t.Errorf("history score %d exceeded MAX_HISTORY (%d)", s.History[piece][target], MAX_HISTORY)
		}

		prePenaltyScore := s.History[piece][target]
		s.updateHistory(move, -500)
		if s.History[piece][target] >= prePenaltyScore {
			t.Errorf("expected penalty to decrease history score from %d, got %d",
				prePenaltyScore, s.History[piece][target])
		}
	})
}
