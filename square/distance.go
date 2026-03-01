package square

// Distance: number of King moves on the otherwise empty board from any square to the four squares {d4, d5, e4, e5} in the center of the board

// Center Chebyshev Distance
// Center Manhattan Distance
var manhattanDist = [64]int8{
	6, 5, 4, 3, 3, 4, 5, 6,
	5, 4, 3, 2, 2, 3, 4, 5,
	4, 3, 2, 1, 1, 2, 3, 4,
	3, 2, 1, 0, 0, 1, 2, 3,
	3, 2, 1, 0, 0, 1, 2, 3,
	4, 3, 2, 1, 1, 2, 3, 4,
	5, 4, 3, 2, 2, 3, 4, 5,
	6, 5, 4, 3, 3, 4, 5, 6,
}

// Calculate Chebyshev Distance
func GetChebyshevDist(square Square) int8 {
	bit1 := uint64(0xFFFFC3C3C3C3FFFF)
	bit0 := uint64(0xFF81BDA5A5BD81FF)

	dist := 2*((bit1>>uint64(square))&1) + ((bit0 >> uint64(square)) & 1)

	return int8(dist)
}

func GetManhattanDist(square Square) int8 {
	file := int8(square & 7)
	rank := int8(square >> 3)

	file ^= (file - 4) >> 8
	rank ^= (rank - 4) >> 8

	return int8(file+rank) & 7
}
