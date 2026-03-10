package main

const (
	BOARD_SQUARE_NUM = 120
)

type Board struct {
	Pieces      []Piece // represents pieces on the 120-index board. the square on which a piece exists has its corresponding value
	SideToMove  Color
	FiftyMove   int
	Ply         int
	HistoryPly  int
	Castling    CastleBit
	Material    []int
	PieceNumber []int // number of each piece on board
	PieceList   []Square
}

func NewBoard() Board {
	board := Board{}

	board.Pieces = make([]Piece, BOARD_SQUARE_NUM)
	board.PieceNumber = make([]int, 13)
	board.PieceList = make([]Square, 140)
	board.SideToMove = White
	board.Material = make([]int, 2)

	return board
}

func (b *Board) PieceIdx(piece Piece) int {
	return (int(piece) * 10) + b.PieceNumber[int(piece)]
}
