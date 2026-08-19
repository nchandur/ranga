package main

import (
	"math/rand/v2"
	"ranga/board"
)

// precomputes attack lookup tables for leaper pieces
func initLeaperAttacks() {
	for square := range 64 {

		sq := board.Square(square)
		// precompute pawn attacks
		board.PawnAttacks[board.White][square] = board.MaskPawnAttacks(sq, board.White)
		board.PawnAttacks[board.Black][square] = board.MaskPawnAttacks(sq, board.Black)

		// precompute knight attacks
		board.KnightAttacks[square] = board.MaskKnightAttacks(sq)

		// precompute king attacks
		board.KingAttacks[square] = board.MaskKingAttacks(sq)
	}

}

// precomputes attack lookup tables for slider pieces
func initSliderAttacks(isBishop bool) {

	for sq := range 64 {
		board.BishopMasks[sq] = board.MaskBishopAttacks(board.Square(sq))
		board.RookMasks[sq] = board.MaskRookAttacks(board.Square(sq))

		attackMask := board.BitBoard(0)

		// set mask and fetch relevant occupancy mask for target piece
		if isBishop {
			attackMask = board.BishopMasks[sq]
		} else {
			attackMask = board.RookMasks[sq]
		}

		relevantBitCount := attackMask.CountBits()
		occupancyIndices := (1 << relevantBitCount)

		// populate magic lookup table for all blocker permutations
		for idx := range occupancyIndices {
			if isBishop {
				occupancy := board.SetOccupancy(idx, relevantBitCount, attackMask)
				magicIdx := (occupancy * board.BishopMagicNumbers[sq]) >> (64 - board.BishopRelevantOccupancy[sq])
				board.BishopAttacks[sq][magicIdx] = board.BishopAttackOTF(board.Square(sq), occupancy)
			} else {
				occupancy := board.SetOccupancy(idx, relevantBitCount, attackMask)
				magicIdx := (occupancy * board.RookMagicNumbers[sq]) >> (64 - board.RookRelevantOccupancy[sq])
				board.RookAttacks[sq][magicIdx] = board.RookAttackOTF(board.Square(sq), occupancy)
			}
		}

	}

}

// populates the Zobrist hashing tables with 64-bit pseudo-random numbers
func initHashKeys() {

	// piece keys
	for pce := board.WP; pce <= board.BK; pce++ {
		for sq := range 64 {
			board.PieceKeys[pce][sq] = rand.Uint64()
		}
	}

	// castle keys
	for pce := range 16 {
		board.CastleKeys[pce] = rand.Uint64()
	}

	// enpassant keys
	for sq := range 64 {
		board.EnpassantKeys[sq] = rand.Uint64()
	}

	// side keys
	board.SideKey = rand.Uint64()

}

func init() {
	initLeaperAttacks()
	initSliderAttacks(true)
	initSliderAttacks(false)
	initHashKeys()
}
