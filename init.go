package main

import "math/rand"

// initializes file and rank boards
func initFileRankBoard() {
	FilesBoard = make([]int, BOARD_SQUARE_NUM)
	RanksBoard = make([]int, BOARD_SQUARE_NUM)

	for idx := range BOARD_SQUARE_NUM {
		FilesBoard[idx] = int(Offboard)
		RanksBoard[idx] = int(Offboard)
	}

	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FR2Square(file, rank)
			FilesBoard[sq] = int(file)
			RanksBoard[sq] = int(rank)
		}
	}

}

// initializes hash slices
func initHashkeys() {
	rng := rand.New(rand.NewSource(19991211))

	for piece := range 13 {
		for sq := range 120 {
			PieceKeys[piece][sq] = rng.Uint64()
		}
	}

	for i := range 16 {
		CastleKeys[i] = rng.Uint64()
	}

	for sq := range 120 {
		EnPassKeys[sq] = rng.Uint64()
	}

	SideKey = rng.Uint64()
}

// initializes board index slices
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

func init() {
	initFileRankBoard()
	initHashkeys()
	initBoardIndices()
}
