package main

const (
	BOARD_SQUARE_NUM = 120
)

var FilesBoard = []int{}
var RanksBoard = []int{}

// returns square index for a given file and rank
func FR2Square(file File, rank Rank) Square {
	return Square((21 + int(file)) + (int(rank) * 10))
}

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
