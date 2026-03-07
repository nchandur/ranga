package board

var fr120To64 = []int{}
var fr64To120 = []int{}

var fr120ToFile = []int{}
var fr120ToRank = []int{}

func InitIndices() {

	fr120To64 = make([]int, 120)
	fr64To120 = make([]int, 64)

	fr120ToFile = make([]int, 120)
	fr120ToRank = make([]int, 120)

	var initIndices = func() {
		sq64 := Square(0)

		for idx := range 120 {
			fr120To64[idx] = 65
		}

		for idx := range 64 {
			fr64To120[idx] = 120
		}

		for rank := One; rank <= Eight; rank++ {
			for file := A; file <= H; file++ {
				sq := FRToSq(file, rank)
				fr64To120[sq64] = int(sq)
				fr120To64[sq] = int(sq64)
				sq64++
			}
		}
	}

	var initFR = func() {

		for idx := range 120 {
			fr120ToFile[idx] = 120
			fr120ToRank[idx] = 120
		}

		for rank := One; rank <= Eight; rank++ {
			for file := A; file <= H; file++ {
				sq := FRToSq(file, rank)
				fr120ToFile[sq] = int(file)
				fr120ToRank[sq] = int(rank)
			}
		}

	}

	initIndices()
	initFR()
}
