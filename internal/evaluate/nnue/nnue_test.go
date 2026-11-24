package nnue

import (
	"ranga/internal/board"
	"testing"
)

// helper to verify that two accumulators match across all hidden dimensions
func assertAccumulatorsEqual(t *testing.T, got, want Accumulator, perspective string) {
	t.Helper()
	for i := range HiddenSize {
		if got.Values[i] != want.Values[i] {
			t.Fatalf("[%s] accumulator mismatch at index %d: got %d, want %d",
				perspective, i, got.Values[i], want.Values[i])
		}
	}
}

// resolves any pending deltas and returns the current ply's accumulators
func current(nn *NNUE) (white, black Accumulator) {
	nn.materialize(int(board.White))
	nn.materialize(int(board.Black))
	e := &nn.stack[nn.ply]
	return e.acc[int(board.White)], e.acc[int(board.Black)]
}

// builds a fresh NNUE sharing nn's network and resets it on b
func scratchFor(nn *NNUE, b *board.Board) *NNUE {
	s := &NNUE{}
	s.Network = nn.Network
	s.Reset(b)
	return s
}

// asserts nn's current (materialized) accumulators match a from-scratch rebuild of b
func assertMatchesScratch(t *testing.T, nn *NNUE, b *board.Board) {
	t.Helper()
	gotW, gotB := current(nn)
	want := scratchFor(nn, b)
	assertAccumulatorsEqual(t, gotW, want.stack[0].acc[int(board.White)], "White")
	assertAccumulatorsEqual(t, gotB, want.stack[0].acc[int(board.Black)], "Black")
}

// finds a legal move by its UCI string
func findMove(t *testing.T, b *board.Board, uci string) board.Move {
	t.Helper()
	var ml board.MoveList
	ml.GenerateMoves(b)
	for i := range ml.Count {
		if ml.Moves[i].String() == uci {
			return ml.Moves[i]
		}
	}
	t.Fatalf("move %s not generated", uci)
	return 0
}

// plays a move on the board and records it in the NNUE stack
func play(t *testing.T, nn *NNUE, b *board.Board, uci string) {
	t.Helper()
	move := findMove(t, b, uci)
	before := b.Preserve()
	if !b.MakeMove(move, false) {
		t.Fatalf("MakeMove returned false for %s", uci)
	}
	nn.Update(&before, move)
}

func TestNNUE(t *testing.T) {
	t.Run("Pop restores exact accumulator state", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		play(t, nn, &b, "e2e4")
		if nn.ply != 1 {
			t.Fatalf("ply after Update = %d; want 1", nn.ply)
		}
		// force the child to be computed, then pop
		current(nn)
		nn.Pop()
		if nn.ply != 0 {
			t.Fatalf("ply after Pop = %d; want 0", nn.ply)
		}

		// the root entry must be untouched: compare against the start position
		start := board.NewBoard()
		start.ParseFEN(board.START)
		assertMatchesScratch(t, nn, &start)
	})

	t.Run("Incremental Update matches scratch Reset across move types", func(t *testing.T) {
		tests := []struct {
			name    string
			fen     string
			uciMove string
		}{
			{"Quiet move e2e4", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2e4"},
			{"Normal capture exd5", "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", "e4d5"},
			{"White Kingside Castling", "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4", "e1g1"},
			{"White Queenside Castling", "r3k2r/pppq1ppp/2npbn2/2b1p3/4P3/2NPBN2/PPPQ1PPP/R3K2R w KQkq - 4 7", "e1c1"},
			{"Black Kingside Castling", "r1bqk2r/pppp1ppp/2n2n2/4p3/1bB1P3/2N2N2/PPPP1PPP/R1BQK2R b KQkq - 5 4", "e8g8"},
			{"Pawn Promotion to Queen", "5k2/4P3/8/8/8/8/8/4K3 w - - 0 1", "e7e8q"},
			{"White En Passant Capture", "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3", "e5f6"},
			{"Black En Passant Capture", "rnbqkbnr/pppp1ppp/8/8/3Pp3/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 2", "e4d3"},
			// king crosses the mirror boundary (e-file -> d-file): forces a refresh
			{"White king mirror crossing", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", "e1d1"},
			{"Black king mirror crossing", "4k3/8/8/8/8/8/8/4K3 b - - 0 1", "e8d8"},
			{"King mirror crossing with pieces", "r3k2r/pppq1ppp/2npbn2/2b1p3/4P3/2NPBN2/PPPQ1PPP/R3K2R w KQkq - 4 7", "e1d1"},
			// king crosses an input bucket boundary (rank 1 -> rank 2)
			{"White king bucket crossing", "4k3/8/8/8/8/8/8/3K4 w - - 0 1", "d1d2"},
			{"King capture-free bucket crossing with pieces", "r3k2r/pppq1ppp/2npbn2/2b1p3/4P3/2NPBN2/PPPQ1PPP/R3K2R w KQkq - 4 7", "e1f1"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				nn := NewRandom()
				b := board.NewBoard()
				b.ParseFEN(tt.fen)
				nn.Reset(&b)

				play(t, nn, &b, tt.uciMove)
				assertMatchesScratch(t, nn, &b)
			})
		}
	})

	t.Run("Lazy delta chain resolves correctly in one Evaluate", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		// several plies with no evaluation in between, including a castle,
		// a capture, and a king move
		for _, m := range []string{
			"e2e4", "e7e5", "g1f3", "b8c6", "f1c4", "g8f6",
			"e1g1", "f6e4", "f1e1", "e4f6", "g1f1",
		} {
			play(t, nn, &b, m)
		}
		if nn.ply != 11 {
			t.Fatalf("ply = %d; want 11", nn.ply)
		}
		assertMatchesScratch(t, nn, &b)
	})

	t.Run("Partially materialized chain stays consistent", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		play(t, nn, &b, "e2e4")
		play(t, nn, &b, "e7e5")
		assertMatchesScratch(t, nn, &b) // materializes plies 1-2

		play(t, nn, &b, "g1f3")
		play(t, nn, &b, "b8c6")
		assertMatchesScratch(t, nn, &b) // builds on the earlier computed entries
	})

	t.Run("Pop then a different move does not reuse stale state", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		// search e2e4, evaluate it, undo, then search d2d4 into the same slot
		before := b.Preserve()
		play(t, nn, &b, "e2e4")
		current(nn)
		nn.Pop()
		b.Restore(&before)

		play(t, nn, &b, "d2d4")
		assertMatchesScratch(t, nn, &b)
	})

	t.Run("PushNull carries the parent accumulators forward", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		play(t, nn, &b, "e2e4")
		parentW, parentB := current(nn)

		nn.PushNull()
		nullW, nullB := current(nn)
		assertAccumulatorsEqual(t, nullW, parentW, "White")
		assertAccumulatorsEqual(t, nullB, parentB, "Black")

		// a real move after the null move still builds on the right base
		nn.Pop()
		nn.Pop()
		if nn.ply != 0 {
			t.Fatalf("ply after two Pops = %d; want 0", nn.ply)
		}
	})

	t.Run("Move after PushNull matches scratch", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		play(t, nn, &b, "e2e4")

		// Black passes (null move), so White is on move again
		b.Side ^= 1
		b.EnPassant = board.NoSquare
		nn.PushNull()

		play(t, nn, &b, "g1f3")
		assertMatchesScratch(t, nn, &b)
	})

	t.Run("outputBucket maps piece count to bucket", func(t *testing.T) {
		tests := []struct {
			name   string
			fen    string
			pieces int
			want   int
		}{
			{"kings only", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", 2, 0},
			{"5 pieces, top of bucket 0", "4k3/p7/8/8/8/8/PP6/4K3 w - - 0 1", 5, 0},
			{"6 pieces, bottom of bucket 1", "4k3/pp6/8/8/8/8/PP6/4K3 w - - 0 1", 6, 1},
			{"10 pieces", "4k3/pppp4/8/8/8/8/PPPP4/4K3 w - - 0 1", 10, 2},
			{"start position", board.START, 32, 7},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b := board.NewBoard()
				b.ParseFEN(tt.fen)
				got := outputBucket(&b)
				if got != tt.want {
					t.Errorf("outputBucket(%d pieces) = %d; want %d", tt.pieces, got, tt.want)
				}
				if got < 0 || got >= OutputBuckets {
					t.Errorf("outputBucket = %d out of range [0, %d)", got, OutputBuckets)
				}
			})
		}
	})

	t.Run("Evaluate alternates us and them based on b.Side and uses the output bucket", func(t *testing.T) {
		fens := []string{
			board.START,
			"4k3/pp6/8/8/8/8/PP6/4K3 w - - 0 1",
			"r3k2r/pppq1ppp/2npbn2/2b1p3/4P3/2NPBN2/PPPQ1PPP/R3K2R w KQkq - 4 7",
			"5k2/4P3/8/8/8/8/8/4K3 w - - 0 1",
		}

		for _, fen := range fens {
			nn := NewRandom()
			b := board.NewBoard()
			b.ParseFEN(fen)
			nn.Reset(&b)
			bucket := outputBucket(&b)
			root := &nn.stack[0]

			b.Side = board.White
			whiteEval := nn.Evaluate(&b)
			expectedWhite := int(nn.Network.Evaluate(&root.acc[int(board.White)], &root.acc[int(board.Black)], bucket))
			if whiteEval != expectedWhite {
				t.Errorf("%s: White eval = %d; want %d", fen, whiteEval, expectedWhite)
			}

			b.Side = board.Black
			blackEval := nn.Evaluate(&b)
			expectedBlack := int(nn.Network.Evaluate(&root.acc[int(board.Black)], &root.acc[int(board.White)], bucket))
			if blackEval != expectedBlack {
				t.Errorf("%s: Black eval = %d; want %d", fen, blackEval, expectedBlack)
			}
		}
	})

	t.Run("Evaluate after Updates matches a fresh Reset", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		for _, m := range []string{"e2e4", "e7e5", "g1f3", "b8c6", "f1c4", "g8f6", "e1g1"} {
			play(t, nn, &b, m)
		}

		got := nn.Evaluate(&b)
		want := scratchFor(nn, &b).Evaluate(&b)
		if got != want {
			t.Errorf("incremental eval = %d; fresh eval = %d", got, want)
		}
	})
}
