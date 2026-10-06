package main

import (
	"fmt"
	"math/rand/v2"
	"ranga/internal/board"
)

func main() {
	pcg := rand.NewPCG(board.SEED, 0)
	rng := rand.New(pcg)

	fmt.Println("var BishopMagic = [64]uint64{")
	for i := range board.Square(64) {
		idx := findIdx(i, board.BishopRelevantOccupancy[i], true, rng)
		fmt.Printf("\t0x%016X,\n", idx)
	}
	fmt.Println("}")

	fmt.Println("var RookMagic = [64]uint64{")
	for i := range board.Square(64) {
		idx := findIdx(i, board.RookRelevantOccupancy[i], false, rng)
		fmt.Printf("\t0x%016X,\n", idx)
	}
	fmt.Println("}")

}

func generateMagicNumber(rng *rand.Rand) uint64 {
	return rng.Uint64() & rng.Uint64() & rng.Uint64()
}

func findIdx(square board.Square, relevantBits int, isBishop bool, rng *rand.Rand) uint64 {
	occ := make([]board.BitBoard, 4096)
	attacks := make([]board.BitBoard, 4096)
	used := make([]board.BitBoard, 4096)

	mask := board.MaskBishopAttacks(square)
	if !isBishop {
		mask = board.MaskRookAttacks(square)
	}

	occIdx := 1 << relevantBits

	for i := range occIdx {
		occ[i] = board.SetOccupancy(i, relevantBits, mask)

		if isBishop {
			attacks[i] = board.BishopAttackOTF(square, occ[i])
		} else {
			attacks[i] = board.RookAttackOTF(square, occ[i])
		}
	}

	for range 100000000 {
		magicNumber := board.BitBoard(generateMagicNumber(rng))

		candidate := (mask * magicNumber) & 0xFF00000000000000
		if candidate.CountBits() < 6 {
			continue
		}

		for i := range 4096 {
			used[i] = 0
		}

		fail := false

		for idx := range occIdx {
			magicIdx := int((occ[idx] * magicNumber) >> (64 - relevantBits))

			if used[magicIdx] == 0 {
				used[magicIdx] = attacks[idx]
			} else if used[magicIdx] != attacks[idx] {
				fail = true
				break
			}
		}

		if !fail {
			return uint64(magicNumber)
		}
	}
	return 0
}
