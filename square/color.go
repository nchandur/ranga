package square

type Color uint8

const (
	LIGHT Color = iota
	DARK
)

// color of the square on 8x8 chessboard
func GetColor(square Square) uint8 {

	// represents all dark squares on board
	// . 1 . 1 . 1 . 1
	// 1 . 1 . 1 . 1 .
	// . 1 . 1 . 1 . 1
	// 1 . 1 . 1 . 1 .
	// . 1 . 1 . 1 . 1
	// 1 . 1 . 1 . 1 .
	// . 1 . 1 . 1 . 1
	// 1 . 1 . 1 . 1 .
	// dark = 1010101001010101101010100101010110101010010101011010101001010101
	dark := uint64(0xAA55AA55AA55AA55)
	return uint8(dark>>square) & 1
}
