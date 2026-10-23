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

func TestNNUE(t *testing.T) {
	t.Run("Preserve and Restore restores exact accumulator state", func(t *testing.T) {
		nn := NewRandom()
		b := board.NewBoard()
		b.ParseFEN(board.START)
		nn.Reset(&b)

		snap := nn.Preserve()
		move := board.NewMove(board.E2, board.E4, board.MoveQuiet)

		nn.Update(&b, move)

		nn.Restore(snap)

		fresh := NewRandom()
		fresh.Network = nn.Network
		fresh.Reset(&b)

		assertAccumulatorsEqual(t, nn.acc.White, fresh.acc.White, "White")
		assertAccumulatorsEqual(t, nn.acc.Black, fresh.acc.Black, "Black")
	})
	t.Run("Incremental Update matches scratch Reset across move types", func(t *testing.T) {
		tests := []struct {
			name    string
			fen     string
			uciMove string
		}{
			{
				name:    "Quiet move e2e4",
				fen:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
				uciMove: "e2e4",
			},
			{
				name:    "Normal capture exd5",
				fen:     "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
				uciMove: "e4d5",
			},
			{
				name:    "White Kingside Castling",
				fen:     "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
				uciMove: "e1g1",
			},
			{
				name:    "White Queenside Castling",
				fen:     "r3k2r/pppq1ppp/2npbn2/2b1p3/4P3/2NPBN2/PPPQ1PPP/R3K2R w KQkq - 4 7",
				uciMove: "e1c1",
			},
			{
				name:    "Black Kingside Castling",
				fen:     "r1bqk2r/pppp1ppp/2n2n2/4p3/1bB1P3/2N2N2/PPPP1PPP/R1BQK2R b KQkq - 5 4",
				uciMove: "e8g8",
			},
			{
				name:    "Pawn Promotion to Queen",
				fen:     "5k2/4P3/8/8/8/8/8/4K3 w - - 0 1",
				uciMove: "e7e8q",
			},
			{
				name:    "White En Passant Capture",
				fen:     "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
				uciMove: "e5f6",
			},
			{
				name:    "Black En Passant Capture",
				fen:     "rnbqkbnr/pppp1ppp/8/8/3Pp3/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 2",
				uciMove: "e4d3",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				nn := NewRandom()
				b := board.NewBoard()
				b.ParseFEN(tt.fen)
				var ml board.MoveList
				ml.GenerateMoves(&b)

				var move board.Move
				found := false
				for i := range ml.Count {
					m := ml.Moves[i]
					if m.String() == tt.uciMove {
						move = m
						found = true
						break
					}
				}

				if !found {
					t.Fatalf("legal move %s not generated for position %s", tt.uciMove, tt.fen)
				}

				nn.Reset(&b)

				copyBoard := b.Preserve()
				if !b.MakeMove(move, false) {
					t.Fatalf("MakeMove returned false for move %v", move)
				}
				nn.Update(&copyBoard, move)

				scratchNN := &NNUE{}
				scratchNN.Network = nn.Network
				scratchNN.Reset(&b)

				assertAccumulatorsEqual(t, nn.acc.White, scratchNN.acc.White, "White")
				assertAccumulatorsEqual(t, nn.acc.Black, scratchNN.acc.Black, "Black")
			})
		}
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

			b.Side = board.White
			whiteEval := nn.Evaluate(&b)
			expectedWhite := int(nn.Network.Evaluate(&nn.acc.White, &nn.acc.Black, bucket))
			if whiteEval != expectedWhite {
				t.Errorf("%s: White eval = %d; want %d", fen, whiteEval, expectedWhite)
			}

			b.Side = board.Black
			blackEval := nn.Evaluate(&b)
			expectedBlack := int(nn.Network.Evaluate(&nn.acc.Black, &nn.acc.White, bucket))
			if blackEval != expectedBlack {
				t.Errorf("%s: Black eval = %d; want %d", fen, blackEval, expectedBlack)
			}
		}
	})
}
