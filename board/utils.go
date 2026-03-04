package board

import "fmt"

// returns 12x10 equivalent index from rank, file
func FRToSq(file File, rank Rank) Square {
	var res int
	// 21st square in the 12x10 corresponds to A1 in the 8x8
	res = int(21+file) + int(rank*10)
	return Square(res)
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
	rank := (idx / 10) - 1

	return fmt.Sprintf("%c%d", rune('a'+file), rank)

}

// returns 64-based index for a given 120-based index
func Fr120To64(idx int) int {

	r := idx / 10
	f := idx % 10

	if f < 1 || f > 8 || r < 2 || r > 9 {
		return 65
	}

	return ((r - 2) * 8) + (f - 1)

}

// return 120-based index for a given 64-based index
func Fr64To120(idx int) int {

	if idx < 0 || idx >= 64 {
		return 120
	}

	r := idx / 8
	f := idx % 8

	return ((r + 2) * 10) + (f + 1)
}
