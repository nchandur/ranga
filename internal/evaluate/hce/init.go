package hce

import (
	"ranga/internal/board"
)

func initEvaluationMasks() {
	var setFileRankMask func(int, int) board.BitBoard = func(rank, file int) board.BitBoard {
		var mask board.BitBoard

		for r := range 8 {
			for f := range 8 {
				square := board.FRtoSq(r, f)

				if file != -1 {
					if f == file {
						mask.SetBit(square)
					}
				} else if rank != -1 {
					if r == rank {
						mask.SetBit(square)
					}
				}

			}
		}

		return mask
	}

	// init file masks
	for rank := range 8 {
		for file := range 8 {
			sq := board.FRtoSq(rank, file)
			FileMasks[sq] |= setFileRankMask(-1, file)
		}
	}

	// init rank masks
	for rank := range 8 {
		for file := range 8 {
			sq := board.FRtoSq(rank, file)
			RankMasks[sq] |= setFileRankMask(rank, -1)
		}
	}

	// init isolated masks
	for rank := range 8 {
		for file := range 8 {
			sq := board.FRtoSq(rank, file)
			IsolatedMasks[sq] |= setFileRankMask(file-1, -1)
			IsolatedMasks[sq] |= setFileRankMask(file+1, -1)
		}
	}

	// init white passed pawn masks

	for rank := range 8 {
		for file := range 8 {
			sq := board.FRtoSq(rank, file)

			WhitePassedPawnMasks[sq] |= setFileRankMask(-1, file-1)
			WhitePassedPawnMasks[sq] |= setFileRankMask(-1, file)
			WhitePassedPawnMasks[sq] |= setFileRankMask(-1, file+1)

			for i := range 8 - rank {
				WhitePassedPawnMasks[sq] &= ^RankMasks[(7-i)*8+file]
			}

		}
	}

	for rank := range 8 {
		for file := range 8 {
			sq := board.FRtoSq(rank, file)

			BlackPassedPawnMasks[sq] |= setFileRankMask(-1, file-1)
			BlackPassedPawnMasks[sq] |= setFileRankMask(-1, file)
			BlackPassedPawnMasks[sq] |= setFileRankMask(-1, file+1)

			for i := range rank + 1 {
				BlackPassedPawnMasks[sq] &= ^RankMasks[i*8+file]
			}

		}
	}

}

func init() {
	initEvaluationMasks()
}
