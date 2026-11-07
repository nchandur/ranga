package nnue

import (
	"ranga/internal/board"
	"testing"
)

// start of the bucket block used by a king on kSq for the given perspective
func bucketBase(perspectiveIsWhite bool, kSq board.Square) int {
	return inputBucket(perspectiveIsWhite, kSq) * FeaturesPerBucket
}

func TestFeatureIndex(t *testing.T) {
	t.Run("indices stay strictly within [0, InputSize) and inside their bucket block", func(t *testing.T) {
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

	t.Run("adds bucket offset for the king square", func(t *testing.T) {
		// same piece and square, kings in different buckets
		for _, isWhite := range []bool{true, false} {
			for kSq := range board.Square(64) {
				mirroredKSq := board.Square(int(kSq) ^ 7)
				if inputBucket(isWhite, kSq) != inputBucket(isWhite, mirroredKSq) {
					t.Errorf("bucket for kSq=%d differs from its mirror %d", kSq, mirroredKSq)
				}
			}
		}

		// a king on white's rank 1 and a king on white's rank 8 must use different blocks
		low := FeatureIndex(true, true, board.WP, board.E4, board.D1)
		high := FeatureIndex(true, true, board.WP, board.E4, board.D8)
		if InputBuckets > 1 && low/FeaturesPerBucket == high/FeaturesPerBucket {
			t.Errorf("expected different bucket blocks for kings on D1 and D8: got %d and %d", low, high)
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

		usIdx := FeatureIndex(false, true, piece, sq, kSq)
		themIdx := FeatureIndex(false, false, piece, sq, kSq)
		base := bucketBase(false, kSq)

		if themIdx-usIdx != 384 {
			t.Errorf("expected 384 offset between us and them: us=%d, them=%d, diff=%d",
				usIdx, themIdx, themIdx-usIdx)
		}
		if usIdx-base >= 384 {
			t.Errorf("friendly piece index %d should be < %d", usIdx, base+384)
		}
		if themIdx-base < 384 {
			t.Errorf("opponent piece index %d should be >= %d", themIdx, base+384)
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
}
