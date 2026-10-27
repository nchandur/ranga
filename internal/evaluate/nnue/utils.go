package nnue

import (
	"math/bits"
	"ranga/internal/board"
)

func screlu(x int16) int32 {
	y := int32(clamp(x, 0, QA))
	return y * y
}

// helper function to clamp bonus between MAX_HISTORY
func clamp(n, low, high int16) int16 {
	if n <= low {
		return low
	}

	if n >= high {
		return high
	}

	return n
}

func lizardDot(v, w []int16) int32 {
	var sum int32
	for i := range v {
		c := clamp(v[i], 0, QA)
		t := c * w[i]
		sum += int32(t) * int32(c)
	}
	return sum
}

func outputBucket(b *board.Board) int {
	count := bits.OnesCount64(uint64(b.Occupancies[2]))
	return (count - 2) / 4
}

// fetch white and black indices for a given piece and square
func getIndices(piece board.Piece, sq board.Square) (int, int) {
	isWhite := int(piece) <= 5
	wIdx := FeatureIndex(true, isWhite, piece, sq)
	bIdx := FeatureIndex(false, !isWhite, piece, sq)
	return wIdx, bIdx
}
