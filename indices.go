package main

var Fr120To64 = []int{}
var Fr64To120 = []int{}

func initBoardIndices() {

	Fr120To64 = make([]int, BOARD_SQUARE_NUM)
	Fr64To120 = make([]int, 64)

	sq64 := 0

	for i := range BOARD_SQUARE_NUM {
		Fr120To64[i] = 65
	}

	for i := range 64 {
		Fr64To120[i] = 120
	}

	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FR2Square(file, rank)
			Fr64To120[sq64] = int(sq)
			Fr120To64[sq] = sq64
			sq64++
		}
	}

}
