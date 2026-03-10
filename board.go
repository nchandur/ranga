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
	EnPassant   Square
	Hash        uint64 // unique value for a given configuration on board
}

func NewBoard() Board {
	board := Board{}

	board.Pieces = make([]Piece, BOARD_SQUARE_NUM)
	board.PieceNumber = make([]int, 13)
	board.PieceList = make([]Square, 130)
	board.SideToMove = White
	board.Material = make([]int, 2)
	board.EnPassant = NoSquare
	board.Hash = 0

	return board
}

func (b *Board) PieceIdx(piece Piece) int {
	return (int(piece) * 10) + b.PieceNumber[int(piece)]
}

func (b *Board) GenHash() uint64 {
	var res uint64

	piece := Empty

	for sq := range BOARD_SQUARE_NUM {
		piece = b.Pieces[sq]

		if piece != Empty && sq != int(Offboard) {
			res ^= PieceKeys[(int(piece)*120)+sq]
		}

	}

	if b.SideToMove == White {
		res ^= SideKey
	}

	if b.EnPassant != NoSquare {
		res ^= PieceKeys[b.EnPassant]
	}

	res ^= CastleKeys[b.Castling]

	return res
}
