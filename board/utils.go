package board

// returns 12x10 equivalent index from rank, file
func FRToSq(file File, rank Rank) Square {
	var res uint8
	// 21st square in the 12x10 corresponds to A1 in the 8x8
	res = uint8(21+file) + uint8(rank*10)
	return Square(res)
}

// assigns index values to 8x8 and 12x10 boards
func Init120to64() {
	Sq120to64 = make([]uint8, 120)
	Sq64to120 = make([]uint8, 64)

	for i := range Sq120to64 {
		Sq120to64[i] = 65
	}

	for i := range Sq64to120 {
		Sq64to120[i] = 120
	}

	var sq64 uint8

	for i := One; i <= Eight; i++ {
		for j := A; j <= H; j++ {
			sq := FRToSq(File(j), Rank(i))
			Sq64to120[sq64] = uint8(sq)
			Sq120to64[sq] = sq64
			sq64++
		}
	}

}
