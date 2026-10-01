package search

import (
	"context"
	"ranga/internal/board"
	"testing"
	"time"
)

// mockEvaluator provides deterministic material balance evaluations
type mockEvaluator struct{}

func (m *mockEvaluator) Evaluate(b *board.Board) int {
	pieceValues := map[board.Piece]int{
		board.WP: 100, board.WN: 320, board.WB: 330, board.WR: 500, board.WQ: 900, board.WK: 20000,
		board.BP: 100, board.BN: 320, board.BB: 330, board.BR: 500, board.BQ: 900, board.BK: 20000,
	}

	whiteScore := 0
	blackScore := 0

	for sq := range 64 {
		p := b.Mailbox[sq]
		if p == board.Empty {
			continue
		}
		if p <= board.WK {
			whiteScore += pieceValues[p]
		} else {
			blackScore += pieceValues[p]
		}
	}

	if b.Side == board.White {
		return whiteScore - blackScore
	}
	return blackScore - whiteScore
}

func TestSearcher(t *testing.T) {
	eval := &mockEvaluator{}

	t.Run("Finds mate in 1", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN("r1bqkb1r/pppp1ppp/2n5/4p3/2B1n3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 4")

		move, score := s.Search(context.Background(), &b, 2)

		wantMove := board.NewMove(board.F3, board.F7, board.WQ, board.Empty, true, false, false, false)
		if move != wantMove {
			t.Errorf("expected mate move %v (Qxf7#), got %v", wantMove, move)
		}
		if score < ISMATE-100 {
			t.Errorf("expected mate score near %d, got %d", ISMATE, score)
		}
	})
	t.Run("Stalemate evaluates to 0", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN("k7/2Q5/1K6/8/8/8/8/8 b - - 0 1")

		score := s.AlphaBeta(context.Background(), &b, -INFINITY, INFINITY, 1)
		if score != 0 {
			t.Errorf("expected stalemate score 0, got %d", score)
		}
	})
	t.Run("Fifty-move rule evaluates to 0", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 100 50")
		b.Ply = 1

		score := s.AlphaBeta(context.Background(), &b, -INFINITY, INFINITY, 1)
		if score != 0 {
			t.Errorf("expected fifty-move draw score 0, got %d", score)
		}
	})
	t.Run("IsRepetition detects repeated position keys", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN(board.START)

		b.Repetition.Table[0] = b.Key
		b.Repetition.Table[1] = 0x111111
		b.Repetition.Table[2] = b.Key
		b.Repetition.Idx = 2
		b.FiftyMove = 4

		if !s.IsRepetition(&b) {
			t.Errorf("expected repetition detected for key %x", b.Key)
		}

		b.Repetition.Table[0] = 0x222222
		b.Repetition.Table[1] = 0x333333
		b.Repetition.Table[2] = b.Key
		b.Repetition.Idx = 2
		b.FiftyMove = 4

		if s.IsRepetition(&b) {
			t.Errorf("did not expect repetition for key %x", b.Key)
		}
	})
	t.Run("Beta cutoff records killer move and history bonus", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN("8/8/8/8/4k3/8/4P3/4K3 w - - 0 1")
		beta := -50000
		alpha := -100000
		depth := 3

		s.AlphaBeta(context.Background(), &b, alpha, beta, depth)

		killer := s.Killers[0][b.Ply]
		if killer == board.NOMOVE {
			t.Errorf("expected killer move recorded on beta cutoff at ply %d", b.Ply)
		}

		if !killer.IsCapture() {
			histScore := s.History[killer.Piece()][killer.Target()]
			if histScore <= 0 {
				t.Errorf("expected positive history score for killer move %v, got %d", killer, histScore)
			}
		}
	})
	t.Run("Search stops promptly when context is cancelled", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN(board.START)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		start := time.Now()
		s.Search(ctx, &b, 20)
		elapsed := time.Since(start)

		if elapsed > 200*time.Millisecond {
			t.Errorf("search took %v; did not terminate promptly on context cancellation", elapsed)
		}
	})
	t.Run("NodeLimit terminates search traversal", func(t *testing.T) {
		s := NewSearcher(eval, 1)
		b := board.NewBoard()
		b.ParseFEN(board.START)

		ctx, cancel := context.WithCancel(context.Background())
		s.Cancel = cancel
		s.NodeLimit = 3000

		s.Search(ctx, &b, 10)

		if s.Nodes > 5000 {
			t.Errorf("search visited %d nodes; expected cutoff near NodeLimit %d", s.Nodes, s.NodeLimit)
		}
	})
}
