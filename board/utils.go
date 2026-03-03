package board

// returns 12x10 equivalent index from rank, file
func FRToSq(file File, rank Rank) Square {
	var res uint8
	// 21st square in the 12x10 corresponds to A1 in the 8x8
	res = uint8(21+file) + uint8(rank*10)
	return Square(res)
}

// returns 64-based index for a given 120-based index
func Fr120To64(idx uint8) uint8 {

	r := idx / 10
	f := idx % 10

	if f < 1 || f > 8 || r < 2 || r > 9 {
		return 65
	}

	return ((r - 2) * 8) + (f - 1)

}

// return 120-based index for a given 64-based index
func Fr64To120(idx uint8) uint8 {

	r := idx / 8
	f := idx % 8

	if r < uint8(A) || r > uint8(H) || f < uint8(One) || f > uint8(Eight) {
		return 120
	}

	return ((r + 2) * 10) + (f + 1)

}
