package nnue

import (
	"ranga/internal/board"
	"testing"
)

func TestFeatureIndex(t *testing.T) {
	t.Run("indices stay strictly within [0, 767] boundary", func(t *testing.T) {
		for _, isWhite := range []bool{true, false} {
			for _, isUs := range []bool{true, false} {
				for pce := board.WP; pce <= board.BK; pce++ {
					for sq := range 64 {
						idx := FeatureIndex(isWhite, isUs, pce, board.Square(sq))
						if idx < 0 || idx >= 768 {
							t.Fatalf("FeatureIndex out of bounds [0, 767]: got %d (perspectiveIsWhite=%v, isUs=%v, piece=%v, sq=%d)",
								idx, isWhite, isUs, pce, sq)
						}
					}
				}
			}
		}
	})
	t.Run("applies rank flip when perspectiveIsWhite is true", func(t *testing.T) {
		sq := board.A1
		wantFlippedSq := int(sq) ^ 56

		idxWhite := FeatureIndex(true, true, board.WP, sq)
		wantWhite := (int(board.WP)%6)*64 + wantFlippedSq
		if idxWhite != wantWhite {
			t.Errorf("white perspective index = %d; want %d", idxWhite, wantWhite)
		}

		idxBlack := FeatureIndex(false, true, board.WP, sq)
		wantBlack := (int(board.WP)%6)*64 + int(sq)
		if idxBlack != wantBlack {
			t.Errorf("black perspective index = %d; want %d", idxBlack, wantBlack)
		}
	})
	t.Run("applies 384 offset for opponent pieces", func(t *testing.T) {
		piece := board.WN
		sq := board.D4

		usIdx := FeatureIndex(false, true, piece, sq)
		themIdx := FeatureIndex(false, false, piece, sq)

		if themIdx-usIdx != 384 {
			t.Errorf("expected 384 offset between us and them: us=%d, them=%d, diff=%d",
				usIdx, themIdx, themIdx-usIdx)
		}
		if usIdx >= 384 {
			t.Errorf("friendly piece index %d should be < 384", usIdx)
		}
		if themIdx < 384 {
			t.Errorf("opponent piece index %d should be >= 384", themIdx)
		}
	})
	t.Run("pieceType modulo partitions pieces across 64-square blocks", func(t *testing.T) {
		sq := board.E4

		for p := range 6 {
			idx := FeatureIndex(false, true, board.Piece(p), sq)
			expectedBlockStart := p * 64
			if idx < expectedBlockStart || idx >= expectedBlockStart+64 {
				t.Errorf("piece %d mapped to index %d; expected range [%d, %d)",
					p, idx, expectedBlockStart, expectedBlockStart+64)
			}
		}
	})
}
