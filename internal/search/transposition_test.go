package search

import (
	"ranga/internal/board"
	"testing"
)

func TestTranspositionTable(t *testing.T) {
	sampleMove := board.NewMove(board.E2, board.E4, board.WP, board.Empty, false, false, false, false)

	t.Run("NewTranspositionTable allocates power-of-two capacity", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		if tt.Length <= 0 {
			t.Fatalf("expected positive table length, got %d", tt.Length)
		}
		if (tt.Length & (tt.Length - 1)) != 0 {
			t.Errorf("Length %d is not a power of two", tt.Length)
		}
		if len(tt.Entries) != tt.Length {
			t.Errorf("slice capacity mismatch: len=%d, Length=%d", len(tt.Entries), tt.Length)
		}

		ttZero := NewTranspositionTable(0)
		if ttZero.Length <= 0 {
			t.Errorf("expected at least 1MB table when size <= 0, got %d", ttZero.Length)
		}
	})
	t.Run("Store and Probe exact score match", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		key := uint64(0x123456789ABCDEF0)

		tt.Store(150, 4, 2, FEXACT, key, sampleMove)

		score := tt.Probe(-1000, 1000, 2, 4, key)
		if score != 150 {
			t.Errorf("Probe score = %d; want 150", score)
		}

		missScore := tt.Probe(-1000, 1000, 2, 5, key)
		if missScore != NOENTRY {
			t.Errorf("expected NOENTRY on deeper probe, got %d", missScore)
		}

		missKey := tt.Probe(-1000, 1000, 2, 4, key^0xFF)
		if missKey != NOENTRY {
			t.Errorf("expected NOENTRY on key miss, got %d", missKey)
		}
	})
	t.Run("Probe alpha and beta cutoffs", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		keyAlpha := uint64(0xAAAA1111)
		keyBeta := uint64(0xBBBB2222)

		tt.Store(50, 3, 0, FALPHA, keyAlpha, sampleMove)

		if got := tt.Probe(100, 200, 0, 3, keyAlpha); got != 100 {
			t.Errorf("FALPHA cutoff expected 100, got %d", got)
		}
		if got := tt.Probe(20, 200, 0, 3, keyAlpha); got != NOENTRY {
			t.Errorf("FALPHA expected NOENTRY when score > alpha, got %d", got)
		}

		tt.Store(250, 3, 0, FBETA, keyBeta, sampleMove)

		if got := tt.Probe(100, 200, 0, 3, keyBeta); got != 200 {
			t.Errorf("FBETA cutoff expected 200, got %d", got)
		}
		if got := tt.Probe(100, 300, 0, 3, keyBeta); got != NOENTRY {
			t.Errorf("FBETA expected NOENTRY when score < beta, got %d", got)
		}
	})
	t.Run("Replacement strategy favors deeper entries", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		key := uint64(0xCAFEBABE)

		tt.Store(200, 5, 0, FEXACT, key, sampleMove)

		shallowMove := board.NewMove(board.D2, board.D4, board.WP, board.Empty, false, false, false, false)
		tt.Store(50, 2, 0, FEXACT, key, shallowMove)

		if got := tt.ProbeMove(key); got != sampleMove {
			t.Errorf("shallow store overwrote deeper entry: got %v, want %v", got, sampleMove)
		}

		deepMove := board.NewMove(board.C2, board.C4, board.WP, board.Empty, false, false, false, false)
		tt.Store(300, 6, 0, FEXACT, key, deepMove)

		if got := tt.ProbeMove(key); got != deepMove {
			t.Errorf("deeper store failed to overwrite: got %v, want %v", got, deepMove)
		}
	})
	t.Run("Mate score normalization across plies", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		key := uint64(0x99998888)

		rawMate := MATESCORE + 10
		storePly := 3
		tt.Store(rawMate, 4, storePly, FEXACT, key, sampleMove)

		idx := key % uint64(tt.Length)
		if tt.Entries[idx].Score != rawMate+storePly {
			t.Errorf("Stored mate score mismatch: got %d, want %d",
				tt.Entries[idx].Score, rawMate+storePly)
		}

		probePly := 1
		gotScore := tt.Probe(-INFINITY, INFINITY, probePly, 4, key)
		wantScore := (rawMate + storePly) - probePly
		if gotScore != wantScore {
			t.Errorf("Normalized probed mate score = %d; want %d", gotScore, wantScore)
		}
	})
	t.Run("ProbeMove and Clear", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		key := uint64(0x5555AAAA)

		tt.Store(100, 3, 0, FEXACT, key, sampleMove)

		if move := tt.ProbeMove(key); move != sampleMove {
			t.Errorf("ProbeMove got %v; want %v", move, sampleMove)
		}
		if move := tt.ProbeMove(key ^ 0x1); move != board.NOMOVE {
			t.Errorf("ProbeMove expected NOMOVE on missing key, got %v", move)
		}

		tt.Clear()

		if move := tt.ProbeMove(key); move != board.NOMOVE {
			t.Errorf("expected NOMOVE after Clear, got %v", move)
		}
		if score := tt.Probe(-1000, 1000, 0, 1, key); score != NOENTRY {
			t.Errorf("expected NOENTRY after Clear, got %d", score)
		}
	})
	t.Run("Resize allocates new capacity", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		oldLength := tt.Length
		tt.Resize(2)

		if tt.Length <= oldLength {
			t.Errorf("expected resized length to be greater than old length, got %d <= %d", tt.Length, oldLength)
		}
		if (tt.Length & (tt.Length - 1)) != 0 {
			t.Errorf("Resized Length %d is not a power of two", tt.Length)
		}
		if len(tt.Entries) != tt.Length {
			t.Errorf("slice capacity mismatch after resize: len=%d, Length=%d", len(tt.Entries), tt.Length)
		}

		key := uint64(0x123456789ABCDEF0)
		missScore := tt.Probe(-1000, 1000, 2, 4, key)
		if missScore != NOENTRY {
			t.Errorf("expected NOENTRY on probe after resize, got %d", missScore)
		}
	})
}
