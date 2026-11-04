package nnue

import (
	"ranga/internal/board"
	"testing"
)

func expectedBucket(white bool, kSq board.Square) int {
	rel := int(kSq)
	if white {
		rel ^= 56
	}
	file := rel & 7
	if file >= 4 {
		file = 7 - file
	}
	return inputBucketLayout[(rel>>3)*4+file]
}

// first feature index of the bucket block used for a king on kSq
func bucketBase(white bool, kSq board.Square) int {
	return expectedBucket(white, kSq) * FeaturesPerBucket
}

func TestFeatureIndex(t *testing.T) {
	t.Run("indices stay within [0, InputSize) and inside the king's bucket block", func(t *testing.T) {
		for _, isWhite := range []bool{true, false} {
			for _, isUs := range []bool{true, false} {
				for pce := board.WP; pce <= board.BK; pce++ {
					for sq := range 64 {
						for kSq := range board.Square(64) {
							idx := FeatureIndex(isWhite, isUs, pce, board.Square(sq), kSq)
							if idx < 0 || idx >= InputSize {
								t.Fatalf("FeatureIndex out of bounds [0, %d): got %d (perspectiveIsWhite=%v, isUs=%v, piece=%v, sq=%d, kSq=%d)",
									InputSize, idx, isWhite, isUs, pce, sq, kSq)
							}
							base := bucketBase(isWhite, kSq)
							if idx < base || idx >= base+FeaturesPerBucket {
								t.Fatalf("FeatureIndex %d outside bucket block [%d, %d) (perspectiveIsWhite=%v, isUs=%v, piece=%v, sq=%d, kSq=%d)",
									idx, base, base+FeaturesPerBucket, isWhite, isUs, pce, sq, kSq)
							}
						}
					}
				}
			}
		}
	})

	t.Run("applies rank flip when perspectiveIsWhite is true", func(t *testing.T) {
		sq := board.A1
		kSq := board.D1

		idxWhite := FeatureIndex(true, true, board.WP, sq, kSq)
		wantFlippedSq := int(sq) ^ 56
		wantWhite := bucketBase(true, kSq) + (int(board.WP)%6)*64 + wantFlippedSq
		if idxWhite != wantWhite {
			t.Errorf("white perspective index = %d; want %d", idxWhite, wantWhite)
		}

		idxBlack := FeatureIndex(false, true, board.WP, sq, kSq)
		wantBlack := bucketBase(false, kSq) + (int(board.WP)%6)*64 + int(sq)
		if idxBlack != wantBlack {
			t.Errorf("black perspective index = %d; want %d", idxBlack, wantBlack)
		}
	})

	t.Run("applies 384 offset for opponent pieces", func(t *testing.T) {
		piece := board.WN
		sq := board.D4
		kSq := board.D1
		base := bucketBase(false, kSq)

		usIdx := FeatureIndex(false, true, piece, sq, kSq)
		themIdx := FeatureIndex(false, false, piece, sq, kSq)

		if themIdx-usIdx != 384 {
			t.Errorf("expected 384 offset between us and them: us=%d, them=%d, diff=%d",
				usIdx, themIdx, themIdx-usIdx)
		}
		if usIdx-base >= 384 {
			t.Errorf("friendly piece index %d should be < base+384 (%d)", usIdx, base+384)
		}
		if themIdx-base < 384 {
			t.Errorf("opponent piece index %d should be >= base+384 (%d)", themIdx, base+384)
		}
	})

	t.Run("pieceType modulo partitions pieces across 64-square blocks", func(t *testing.T) {
		sq := board.E4
		kSq := board.D1
		base := bucketBase(false, kSq)

		for p := range 6 {
			idx := FeatureIndex(false, true, board.Piece(p), sq, kSq)
			expectedBlockStart := base + p*64
			if idx < expectedBlockStart || idx >= expectedBlockStart+64 {
				t.Errorf("piece %d mapped to index %d; expected range [%d, %d)",
					p, idx, expectedBlockStart, expectedBlockStart+64)
			}
		}
	})

	t.Run("applies horizontal mirroring when king is on kingside", func(t *testing.T) {
		piece := board.WP

		idxKingside := FeatureIndex(false, true, piece, board.H1, board.E1)
		idxQueenside := FeatureIndex(false, true, piece, board.A1, board.D1)

		if idxKingside != idxQueenside {
			t.Errorf("expected horizontal mirroring to map H1 (kingside king) to A1 (queenside king): got %d and %d",
				idxKingside, idxQueenside)
		}
	})

	t.Run("mirroring is consistent for every king square, perspective and piece square", func(t *testing.T) {
		for _, isWhite := range []bool{true, false} {
			for kSq := range board.Square(64) {
				if int(kSq)&7 < 4 {
					continue
				}
				mirrorK := board.Square(int(kSq) ^ 7)
				for sq := range 64 {
					a := FeatureIndex(isWhite, true, board.WQ, board.Square(sq), kSq)
					b := FeatureIndex(isWhite, true, board.WQ, board.Square(sq^7), mirrorK)
					if a != b {
						t.Fatalf("mirror mismatch (white=%v, kSq=%d, sq=%d): %d != %d", isWhite, kSq, sq, a, b)
					}
				}
			}
		}
	})

	t.Run("every input bucket is reachable from both perspectives", func(t *testing.T) {
		for _, isWhite := range []bool{true, false} {
			seen := make([]bool, InputBuckets)
			for kSq := range board.Square(64) {
				idx := FeatureIndex(isWhite, true, board.WP, board.A1, kSq)
				seen[idx/FeaturesPerBucket] = true
			}
			for b, ok := range seen {
				if !ok {
					t.Errorf("bucket %d unreachable for perspectiveIsWhite=%v", b, isWhite)
				}
			}
		}
	})
}

func TestInputBucket(t *testing.T) {
	t.Run("lookup table matches the layout for every square", func(t *testing.T) {
		for _, white := range []bool{true, false} {
			for kSq := range board.Square(64) {
				if got, want := inputBucket(white, kSq), expectedBucket(white, kSq); got != want {
					t.Errorf("inputBucket(white=%v, kSq=%d) = %d; want %d", white, kSq, got, want)
				}
			}
		}
	})

	t.Run("buckets are in range and symmetric", func(t *testing.T) {
		for kSq := range board.Square(64) {
			for _, white := range []bool{true, false} {
				b := inputBucket(white, kSq)
				if b < 0 || b >= InputBuckets {
					t.Errorf("inputBucket(white=%v, kSq=%d) = %d out of range", white, kSq, b)
				}
				// file mirror
				if m := inputBucket(white, board.Square(int(kSq)^7)); m != b {
					t.Errorf("file-mirrored king bucket differs (white=%v, kSq=%d): %d vs %d", white, kSq, b, m)
				}
			}
			// colour symmetry: white king on sq sees what a black king on sq^56 sees
			if w, bl := inputBucket(true, kSq), inputBucket(false, board.Square(int(kSq)^56)); w != bl {
				t.Errorf("colour symmetry broken at kSq=%d: white=%d black(sq^56)=%d", kSq, w, bl)
			}
		}
	})

	t.Run("4-bucket layout: back rank / rank 2 / ranks 3-4 / ranks 5-8", func(t *testing.T) {
		if InputBuckets != 4 {
			t.Skipf("hard-coded expectations are for 4 buckets, have %d", InputBuckets)
		}
		tests := []struct {
			name  string
			white bool
			kSq   board.Square
			want  int
		}{
			{"white e1", true, board.E1, 0},
			{"white g1 (mirrored)", true, board.G1, 0},
			{"white e2", true, board.E2, 1},
			{"white h2 (mirrored)", true, board.H2, 1},
			{"white e3", true, board.E3, 2},
			{"white e4", true, board.E4, 2},
			{"white e5", true, board.E5, 3},
			{"white e8", true, board.E8, 3},
			{"black e8", false, board.E8, 0},
			{"black e7", false, board.E7, 1},
			{"black e6", false, board.E6, 2},
			{"black e5", false, board.E5, 2},
			{"black e4", false, board.E4, 3},
			{"black e1", false, board.E1, 3},
		}
		for _, tc := range tests {
			if got := inputBucket(tc.white, tc.kSq); got != tc.want {
				t.Errorf("%s: bucket = %d; want %d", tc.name, got, tc.want)
			}
		}
	})
}

func TestNeedsRefresh(t *testing.T) {
	if InputBuckets != 4 {
		t.Skipf("hard-coded expectations are for 4 buckets, have %d", InputBuckets)
	}
	tests := []struct {
		name     string
		white    bool
		from, to board.Square
		want     bool
	}{
		{"same bucket, no mirror change (d1-c1)", true, board.D1, board.C1, false},
		{"same bucket, both mirrored (e1-f1)", true, board.E1, board.F1, false},
		{"mirror boundary crossed (e1-d1)", true, board.E1, board.D1, true},
		{"bucket 0 -> 1 (e1-e2)", true, board.E1, board.E2, true},
		{"within bucket 2 (d3-d4)", true, board.D3, board.D4, false},
		{"bucket 2 -> 3 (d4-d5)", true, board.D4, board.D5, true},
		{"black bucket 0 -> 1 (e8-e7)", false, board.E8, board.E7, true},
		{"black within bucket 0 (d8-c8)", false, board.D8, board.C8, false},
	}
	for _, tc := range tests {
		if got := needsRefresh(tc.white, tc.from, tc.to); got != tc.want {
			t.Errorf("%s: needsRefresh = %v; want %v", tc.name, got, tc.want)
		}
	}
}
