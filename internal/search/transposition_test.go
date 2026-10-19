package search

import (
	"ranga/internal/board"
	"sync"
	"testing"
)

func TestTranspositionTable(t *testing.T) {
	sampleMove := board.NewMove(board.E2, board.E4, board.MoveQuiet)

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

		shallowMove := board.NewMove(board.D2, board.D4, board.MoveQuiet)
		tt.Store(50, 2, 0, FEXACT, key, shallowMove)

		if got := tt.ProbeMove(key); got != sampleMove {
			t.Errorf("shallow store overwrote deeper entry: got %v, want %v", got, sampleMove)
		}

		deepMove := board.NewMove(board.C2, board.C4, board.MoveQuiet)
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

		stored := TranspositionTableEntry(tt.slot(key).Load())
		if stored.Score() != rawMate+storePly {
			t.Errorf("Stored mate score mismatch: got %d, want %d",
				stored.Score(), rawMate+storePly)
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
	t.Run("Negative mate score normalization across plies", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		key := uint64(0x123456789ABCDEF0)

		rawMated := -MATESCORE - 10
		storePly := 3
		tt.Store(rawMated, 4, storePly, FEXACT, key, sampleMove)

		stored := TranspositionTableEntry(tt.slot(key).Load())
		if stored.Score() != rawMated-storePly {
			t.Errorf("Stored mated score mismatch: got %d, want %d", stored.Score(), rawMated-storePly)
		}

		probePly := 1
		gotScore := tt.Probe(-INFINITY, INFINITY, probePly, 4, key)
		wantScore := (rawMated - storePly) + probePly
		if gotScore != wantScore {
			t.Errorf("Normalized probed mated score = %d; want %d", gotScore, wantScore)
		}
	})
	t.Run("Packed entry round-trips every field", func(t *testing.T) {
		mv := board.NewMove(board.E7, board.E8, board.MovePromoCaptureQ)
		key := uint64(0xDEADBEEFCAFEF00D)

		scores := []int{0, 1, -1, 12345, -12345, INFINITY, -INFINITY, MATESCORE + MAX_PLY, -MATESCORE - MAX_PLY}
		depths := []int{0, 1, 17, 255}
		flags := []int{FEXACT, FALPHA, FBETA}

		for _, score := range scores {
			for _, depth := range depths {
				for _, flag := range flags {
					e := packEntry(key, score, depth, flag, mv)

					if e == 0 {
						t.Fatalf("valid entry packed to 0 (score %d depth %d flag %d)", score, depth, flag)
					}
					if e.Score() != score || e.Depth() != depth || e.Flag() != flag || e.Move() != mv {
						t.Fatalf("fields differ: got score %d depth %d flag %d move %v; want %d %d %d %v",
							e.Score(), e.Depth(), e.Flag(), e.Move(), score, depth, flag, mv)
					}
					if !e.matches(key) {
						t.Fatalf("entry does not match its own key")
					}
				}
			}
		}
	})
	t.Run("Packed fields do not bleed into each other", func(t *testing.T) {
		deep := packEntry(0, 0, 255, FEXACT, board.NOMOVE)
		if deep.Score() != 0 || deep.Flag() != FEXACT || deep.Move() != board.NOMOVE || deep.Depth() != 255 {
			t.Error("depth bled into another field")
		}

		neg := packEntry(0, -1, 0, FBETA, board.NewMove(board.H1, board.A8, board.MovePromoCaptureQ))
		if neg.Depth() != 0 || neg.Score() != -1 || neg.Flag() != FBETA {
			t.Error("score or move bled into another field")
		}
	})
	t.Run("Score range fits the score field", func(t *testing.T) {
		if INFINITY+MAX_PLY >= scoreBias-1 {
			t.Fatalf("scores up to %d do not fit in %d score bits", INFINITY+MAX_PLY, scoreBits)
		}
	})
	t.Run("Key check uses the high key bits only", func(t *testing.T) {
		key := uint64(0xDEADBEEFCAFEF00D)
		e := packEntry(key, 5, 5, FEXACT, board.NOMOVE)

		if !e.matches(key ^ 1) { // low bits only select the slot
			t.Error("low key bits must not affect the check")
		}
		if e.matches(key ^ (1 << 63)) {
			t.Error("different top key bits must not match")
		}
		if TranspositionTableEntry(0).matches(0) {
			t.Error("empty slot must never match")
		}
	})
	t.Run("Capacity uses the full memory budget", func(t *testing.T) {
		for _, mb := range []int{1, 2, 16, 100, 256} {
			tt := NewTranspositionTable(mb)

			bytes := tt.Length * 8
			if bytes > mb*1024*1024 || bytes*2 <= mb*1024*1024 {
				t.Errorf("%dMB: %d entries use %d bytes", mb, tt.Length, bytes)
			}
		}

		tt := NewTranspositionTable(1)
		tt.Resize(16)
		if tt.Length != 1<<21 {
			t.Errorf("16MB should hold 2^21 entries, got %d", tt.Length)
		}
	})
	t.Run("Concurrent Store and Probe stay consistent", func(t *testing.T) {
		tt := NewTranspositionTable(1)
		mv := board.NewMove(board.G1, board.F3, board.MoveQuiet)

		var wg sync.WaitGroup
		for g := range 8 {
			wg.Add(1)
			go func(seed uint64) {
				defer wg.Done()
				x := seed*0x9E3779B97F4A7C15 + 1

				for range 200000 {
					x ^= x << 13
					x ^= x >> 7
					x ^= x << 17

					depth := int(x>>8) & 31
					tt.Store(depth*7, depth, 0, FEXACT, x, mv)

					e := TranspositionTableEntry(tt.Entries[(x>>20)&uint64(tt.Length-1)].Load())
					if e != 0 && (e.Score() != e.Depth()*7 || e.Move() != mv) {
						t.Errorf("inconsistent entry: score %d depth %d move %v", e.Score(), e.Depth(), e.Move())
						return
					}
				}
			}(uint64(g + 1))
		}
		wg.Wait()
	})
}
