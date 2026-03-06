package board

import "fmt"

// from file and rank to 120-based index
func FRToSq(file File, rank Rank) Square {
	var res int
	// 21st square in the 12x10 corresponds to A1 in the 8x8
	res = int(21+file) + int(rank*10)
	return Square(res)
}

// from 120-based index to file
func Fr120ToFile(idx int) File {
	file := (idx % 10)

	return File(file)
}

// from 120-based index to rank
func Fr120ToRank(idx int) Rank {
	rank := (idx / 10)
	return Rank(rank)
}

// from 120-based to 64-based
func Fr120To64(idx int) int {

	r := Fr120ToRank(idx)
	f := Fr120ToFile(idx)

	if f < 1 || f > 8 || r < 2 || r > 9 {
		return 65
	}

	return ((int(r) - 2) * 8) + (int(f) - 1)

}

// from 64-based to 120-based
func Fr64To120(idx int) int {

	if idx < 0 || idx >= 64 {
		return 120
	}

	r := idx / 8
	f := idx % 8

	return ((r + 2) * 10) + (f + 1)
}

// returns file and rank from 64-based index
func Fr64ToFR(idx int) string {

	file := idx % 8
	rank := idx / 8

	return fmt.Sprintf("%c%d", rune('a'+file), rank)

}

// returns file and tank from 120-based index
func Fr120ToFR(idx int) string {

	file := (idx % 10) - 1
	rank := (idx / 10) - 2

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return "-"
	}

	return fmt.Sprintf("%c%d", rune('a'+file), rank+1)

}
