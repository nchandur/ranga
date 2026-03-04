package board

import "math/rand"

var PieceKeys [13][120]uint64
var SideKey uint64
var CastleKeys [16]uint64

// generates a unique 64-bit random number
func rgn64() uint64 {
	var res uint64
	res = rand.Uint64() | (rand.Uint64() << 15) | (rand.Uint64() << 30) | (rand.Uint64() << 45) | (rand.Uint64()&0xf)<<60
	return res
}

func InitHashKeys() {

	for i := range 13 {
		for j := range 120 {
			PieceKeys[i][j] = rgn64()
		}
	}

	SideKey = rgn64()

	for i := range 16 {
		CastleKeys[i] = rgn64()
	}

}
