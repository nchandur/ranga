package uci

import (
	"math/rand/v2"
	"strings"
	"testing"

	"ranga/internal/board"
)

const (
	kiwipete    = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	stalemateFE = "7k/5Q2/6K1/8/8/8/8/8 b - - 0 1"
	checkmateFE = "R5k1/5ppp/8/8/8/8/8/6K1 b - - 0 1"
	forcedFEN   = "k7/8/8/8/8/8/1q6/K7 w - - 0 1"
)

func mustParse(t *testing.T, fen string) board.Board {
	t.Helper()
	b := board.NewBoard()
	if err := b.ParseFEN(fen); err != nil {
		t.Fatalf("bad test FEN %q: %v", fen, err)
	}
	return b
}

func newRNG(seed uint64) *rand.Rand { return rand.New(rand.NewPCG(seed, seed)) }

// position identity ignoring move counters
func key4(b *board.Board) string {
	return strings.Join(strings.Fields(b.FEN())[:4], " ")
}

// every position reachable by one legal move from b
func legalSuccessors(b *board.Board) map[string]struct{} {
	set := map[string]struct{}{}
	list := board.NewMoveList()
	list.GenerateMoves(b)
	for i := range list.Count {
		c := b.Preserve()
		if c.MakeMove(list.Moves[i], false) {
			set[key4(&c)] = struct{}{}
		}
	}
	return set
}

func TestHasLegalMove(t *testing.T) {
	cases := []struct {
		name string
		fen  string
		want bool
	}{
		{"startpos", board.START, true},
		{"in check, king can step away", "7k/8/8/8/8/8/8/R3K2r w - - 0 1", true},
		{"in check, only move is a capture", forcedFEN, true},
		{"stalemate", stalemateFE, false},
		{"checkmate", checkmateFE, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := mustParse(t, tc.fen)
			before := b.Preserve()

			if got := hasLegalMove(&b); got != tc.want {
				t.Errorf("hasLegalMove = %v, want %v", got, tc.want)
			}
			if b != before {
				t.Error("hasLegalMove modified the board")
			}
		})
	}
}

func TestPlayRandomMoveIsLegalAndUnbiased(t *testing.T) {
	cases := []struct {
		name  string
		fen   string
		moves int
	}{
		{"startpos", board.START, 20},
		{"kiwipete", kiwipete, 48},
		{"en passant available", "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3", 0},
		{"castling", "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", 0},
		{"promotion", "8/1P4k1/8/8/8/8/4K1p1/8 w - - 0 1", 0},
		{"black to move", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			start := mustParse(t, tc.fen)
			legal := legalSuccessors(&start)
			if tc.moves != 0 && len(legal) != tc.moves {
				t.Fatalf("test oracle found %d legal moves, want %d: movegen/MakeMove is broken", len(legal), tc.moves)
			}

			rng := newRNG(1)
			seen := map[string]struct{}{}
			for range 5000 {
				b := start.Preserve()
				if !playRandomMove(&b, rng) {
					t.Fatal("playRandomMove returned false on a position with legal moves")
				}
				key := key4(&b)
				if _, ok := legal[key]; !ok {
					t.Fatalf("played an illegal move, reaching %s", b.FEN())
				}
				if b.Ply != start.Ply+1 {
					t.Fatalf("ply = %d, want %d", b.Ply, start.Ply+1)
				}
				if b.Side != start.Side^1 {
					t.Fatalf("side to move did not flip")
				}
				seen[key] = struct{}{}
			}
			if len(seen) != len(legal) {
				t.Errorf("only %d of %d legal moves were ever chosen", len(seen), len(legal))
			}
		})
	}
}

func TestPlayRandomMoveFindsOnlyLegalMove(t *testing.T) {
	start := mustParse(t, forcedFEN)
	want := mustParse(t, "k7/8/8/8/8/8/1K6/8 b - - 0 1")

	for seed := range uint64(50) {
		b := start.Preserve()
		if !playRandomMove(&b, newRNG(seed)) {
			t.Fatalf("seed %d: returned false although Kxb2 is legal", seed)
		}
		if key4(&b) != key4(&want) {
			t.Fatalf("seed %d: got %s, want Kxb2", seed, b.FEN())
		}
	}
}

func TestPlayRandomMoveWithNoLegalMoves(t *testing.T) {
	for name, fen := range map[string]string{"stalemate": stalemateFE, "checkmate": checkmateFE} {
		t.Run(name, func(t *testing.T) {
			start := mustParse(t, fen)
			for seed := range uint64(20) {
				b := start.Preserve()
				if playRandomMove(&b, newRNG(seed)) {
					t.Fatalf("seed %d: returned true on a terminal position", seed)
				}
				if b != start {
					t.Fatalf("seed %d: board was modified by a failed playRandomMove", seed)
				}
			}
		})
	}
}
