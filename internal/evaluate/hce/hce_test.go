package hce

import (
	"ranga/internal/board"
	"testing"
)

func TestHCE_Evaluate(t *testing.T) {
	h := HCE{}

	t.Run("symmetrical starting position evaluates close to zero", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN(board.START)

		score := h.Evaluate(&b)
		if score < -10 || score > 10 {
			t.Errorf("expected near-zero score for starting position, got %d", score)
		}
	})
	t.Run("side to move flips score perspective", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("4k3/8/8/8/8/8/4R3/4K3 w - - 0 1")
		whiteScore := h.Evaluate(&b)

		if whiteScore <= 0 {
			t.Fatalf("expected positive score for White with extra Rook, got %d", whiteScore)
		}

		b.Side = board.Black
		blackScore := h.Evaluate(&b)

		if blackScore != -whiteScore {
			t.Errorf("perspective flip failed: whiteScore=%d, blackScore=%d (expected %d)",
				whiteScore, blackScore, -whiteScore)
		}
	})
	t.Run("material advantage dominates evaluation", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("4k3/8/8/8/8/8/4Q3/4K3 w - - 0 1")
		queenAdvScore := h.Evaluate(&b)
		b.ParseFEN("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1")
		pawnAdvScore := h.Evaluate(&b)

		if queenAdvScore <= pawnAdvScore {
			t.Errorf("queen advantage score (%d) should be much higher than pawn advantage score (%d)",
				queenAdvScore, pawnAdvScore)
		}

		if queenAdvScore < 500 {
			t.Errorf("expected queen advantage score >= 500, got %d", queenAdvScore)
		}
	})
	t.Run("board mirror symmetry invariant", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("r1bqk2r/pp2bppp/2n1pn2/2pp4/3P4/2NBPN2/PPP2PPP/R1BQK2R w KQkq - 2 7")
		whiteScore := h.Evaluate(&b)

		b.Mirror()
		blackScore := h.Evaluate(&b)

		if whiteScore != blackScore {
			t.Errorf("evaluation mirror symmetry broken: original=%d, mirrored=%d",
				whiteScore, blackScore)
		}
	})
	t.Run("empty board does not panic and returns zero", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("8/8/8/8/8/8/8/8 w - - 0 1")

		score := h.Evaluate(&b)
		if score != 0 {
			t.Errorf("expected 0 for empty board, got %d", score)
		}
	})
}
