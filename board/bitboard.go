package board

import (
	"fmt"
)

type BitBoard uint64

// pops the least significant bit from bitboard
func (b *BitBoard) PopBit() int {

	// no pieces on board
	if *b == 0 {
		return -1
	}

	// find the lsb
	lsb := *b & -(*b)
	sq := 0

	for lsb > 1 {
		lsb >>= 1
		sq++
	}
	*b &= *b - 1
	return sq
}

// counts the number of bits in the bitboard aka number of pieces in the bitboard
func (b BitBoard) CountBits() int {
	var res int

	for uint64(b) != 0 {
		res++
		b &= b - 1
	}

	return res
}

// clear bit in square
func (b *BitBoard) ClearBit(square Square) {
	idx := Fr120To64(int(square))

	*b &= BitBoard(ClearMask[idx])
}

// set bitin square
func (b *BitBoard) SetBit(square Square) {
	idx := Fr120To64(int(square))
	*b |= BitBoard(SetMask[idx])
}

// prints all the squares on which pieces exist
func (b BitBoard) Print() {
	shifter := uint64(1)

	for i := One; i <= Eight; i++ {
		for j := A; j <= H; j++ {
			sq := FRToSq(File(j), Rank(i))
			sq64 := Fr120To64(int(sq))

			if (shifter<<uint64(sq64))&uint64(b) != 0 {
				fmt.Printf("%5s", "X")
			} else {
				fmt.Printf("%5s", "-")
			}

		}
		fmt.Println()
	}

}
