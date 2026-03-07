package board

// from file and rank to 120-based index
func FRToSq(file File, rank Rank) Square {
	var res int
	// 21st square in the 12x10 corresponds to A1 in the 8x8
	res = int(21+file) + int(rank*10)
	return Square(res)
}

// from 120-based to 64-based
func Fr120To64(idx int) int {
	return fr120To64[idx]
}

// from 64-based to 120-based
func Fr64To120(idx int) int {
	return fr64To120[idx]
}

// from 120-based index to file
func Fr120ToFile(idx int) File {
	return File(fr120ToFile[idx])
}

// from 120-based index to rank
func Fr120ToRank(idx int) Rank {
	return Rank(fr120ToRank[idx])
}
