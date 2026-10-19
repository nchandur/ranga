package search

import (
	"math/rand"
	"slices"
	"testing"

	"ranga/internal/board"
)

func setup(game, path []uint64, fifty int) (*Searcher, *board.Board) {
	s := &Searcher{}
	s.SetGame(game)
	copy(s.path[:], path)

	b := &board.Board{}
	b.Ply = len(path) - 1
	b.Key = path[len(path)-1]
	b.FiftyMove = fifty

	return s, b
}

func TestIsRepetition(t *testing.T) {
	tests := []struct {
		name  string
		game  []uint64
		path  []uint64 // last entry is the current node
		fifty int
		want  bool
	}{
		{"empty history, root only", nil, []uint64{1}, 0, false},
		{"root only with fifty counter", nil, []uint64{1}, 10, false},
		{"repeat inside game", []uint64{1, 2, 3}, []uint64{1}, 3, true},
		{"no repeat inside game", []uint64{1, 2, 3}, []uint64{4}, 3, false},
		{"repeat outside fifty window", []uint64{1, 2, 3}, []uint64{1}, 2, false},
		{"repeat at edge of window", []uint64{1, 2, 3}, []uint64{1}, 3, true},
		{"repeat on the path", nil, []uint64{1, 2, 3, 1}, 3, true},
		{"path repeat outside window", nil, []uint64{1, 2, 3, 1}, 2, false},
		{"window spans game and path", []uint64{9, 8}, []uint64{5, 6, 9}, 4, true},
		{"window spans game and path, too short", []uint64{9, 8}, []uint64{5, 6, 9}, 3, false},
		{"fifty zero right after irreversible move", []uint64{1, 2}, []uint64{3, 1}, 0, false},
		{"fifty larger than history", []uint64{7}, []uint64{1, 7}, 100, true},
		{"fifty larger than history, no repeat", []uint64{7}, []uint64{1, 2}, 100, false},
		{"current node is not compared with itself", nil, []uint64{1, 2}, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, b := setup(tt.game, tt.path, tt.fifty)
			if got := s.IsRepetition(b); got != tt.want {
				t.Fatalf("IsRepetition = %v, want %v", got, tt.want)
			}
		})
	}
}

// stale path entries above current ply must never be read
func TestIsRepetitionIgnoresStalePath(t *testing.T) {
	s, b := setup(nil, []uint64{1, 2}, 5)
	s.path[2] = 2
	s.path[3] = 2

	if s.IsRepetition(b) {
		t.Fatal("stale entries above the current ply were read")
	}
}

func TestGameHistory(t *testing.T) {
	s := &Searcher{}

	s.AddGameMove(10, 1)
	s.AddGameMove(20, 2)
	if want := []uint64{10, 20}; !slices.Equal(s.game, want) {
		t.Fatalf("game = %v, want %v", s.game, want)
	}

	// irreversible move drops everything before it
	s.AddGameMove(30, 0)
	if len(s.game) != 0 {
		t.Fatalf("game = %v after irreversible move, want empty", s.game)
	}

	s.AddGameMove(40, 1)
	s.ClearGame()
	if len(s.game) != 0 {
		t.Fatalf("game = %v after ClearGame, want empty", s.game)
	}

	keys := []uint64{1, 2, 3}
	s.SetGame(keys)
	keys[0] = 99
	if want := []uint64{1, 2, 3}; !slices.Equal(s.game, want) {
		t.Fatalf("game = %v, want %v (SetGame must copy)", s.game, want)
	}

	s.SetGame([]uint64{5})
	if want := []uint64{5}; !slices.Equal(s.game, want) {
		t.Fatalf("game = %v, want %v", s.game, want)
	}
}

func referenceIsRepetition(table []uint64, idx, fifty int) bool {
	for i := max(idx-fifty, 0); i < idx; i++ {
		if table[i] == table[idx] {
			return true
		}
	}
	return false
}

func TestIsRepetitionMatchesFlatTable(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	for iter := range 20000 {
		gameLen := rng.Intn(12)
		pathLen := 1 + rng.Intn(12)
		fifty := rng.Intn(30)

		flat := make([]uint64, gameLen+pathLen)
		for i := range flat {
			flat[i] = uint64(rng.Intn(6))
		}

		s, b := setup(flat[:gameLen], flat[gameLen:], fifty)

		want := referenceIsRepetition(flat, len(flat)-1, fifty)
		if got := s.IsRepetition(b); got != want {
			t.Fatalf("iter %d: flat=%v game=%d fifty=%d: got %v, want %v",
				iter, flat, gameLen, fifty, got, want)
		}
	}
}

func TestTrimmingLosesNoRepetitions(t *testing.T) {
	rng := rand.New(rand.NewSource(2))

	for iter := range 2000 {
		s := &Searcher{}
		var full []uint64
		fifty := rng.Intn(5)
		key := uint64(rng.Intn(5))
		full = append(full, key)

		for move := range 60 {
			prev := key
			key = uint64(rng.Intn(5))
			if rng.Intn(8) == 0 {
				fifty = 0
			} else {
				fifty++
			}
			full = append(full, key)
			s.AddGameMove(prev, fifty)

			s.path[0] = key
			b := &board.Board{Key: key, FiftyMove: fifty, Ply: 0}

			want := referenceIsRepetition(full, len(full)-1, fifty)
			if got := s.IsRepetition(b); got != want {
				t.Fatalf("iter %d move %d: got %v, want %v (full=%v fifty=%d)",
					iter, move, got, want, full, fifty)
			}
		}
	}
}

// the game history stays small however long the game is
func TestGameHistoryIsBounded(t *testing.T) {
	s := &Searcher{}
	fifty := 0
	for i := range 10000 {
		fifty++
		if fifty == 100 {
			fifty = 0
		}
		s.AddGameMove(uint64(i), fifty)
		if len(s.game) > 100 {
			t.Fatalf("game grew to %d entries", len(s.game))
		}
	}
}
