package main

import (
	"fmt"
	"strings"
)

const (
	BOARD_SQUARE_NUM   = 120
	MAX_GAME_MOVES     = 4096
	MAX_POSITION_MOVES = 256
	MAX_DEPTH          = 64
)

const (
	START     = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	PieceChar = ".PNBRQKpnbqrk"
	SideChar  = "wb-"
	RankChar  = "12345678"
	FileChar  = "abcdefgh"
)

type Board struct {
	Pieces        []Piece // represents pieces on the 120-index board. the square on which a piece exists has its corresponding value
	SideToMove    Color
	FiftyMove     int
	Ply           int
	HistoryPly    int
	Castling      CastleBit
	Material      []int
	PieceNumber   []int // number of each piece on board
	PieceList     []Square
	EnPassant     Square
	Hash          uint64 // unique value for a given configuration on board
	MoveList      []int
	MoveScores    []int
	MoveListStart []int
}

func NewBoard() Board {
	board := Board{}

	board.Pieces = make([]Piece, BOARD_SQUARE_NUM)
	board.PieceNumber = make([]int, 13)
	board.PieceList = make([]Square, 130)
	board.SideToMove = Both
	board.Material = make([]int, 2)
	board.EnPassant = NoSquare
	board.Hash = 0

	board.MoveList = make([]int, MAX_DEPTH*MAX_POSITION_MOVES)
	board.MoveScores = make([]int, MAX_DEPTH*MAX_POSITION_MOVES)
	board.MoveListStart = make([]int, MAX_DEPTH)

	return board
}

func (b *Board) PieceIdx(piece Piece) int {
	return (int(piece) * 10) + b.PieceNumber[int(piece)]
}

func (b *Board) GenHash() uint64 {

	var key uint64

	// pieces
	for sq := range 120 {
		piece := b.Pieces[sq]
		if piece != Empty && piece != 120 {
			key ^= PieceKeys[piece][sq]
		}
	}

	// side to move
	if b.SideToMove == White {
		key ^= SideKey
	}

	// en passant
	if b.EnPassant != NoSquare {
		key ^= EnPassKeys[b.EnPassant]
	}

	// castling rights
	key ^= CastleKeys[b.Castling]

	return key

}
func (b *Board) Reset() {

	for idx := range BOARD_SQUARE_NUM {
		b.Pieces[idx] = 120
	}

	for idx := range 64 {
		b.Pieces[Fr64To120[idx]] = Empty
	}

	for idx := range b.PieceList {
		b.PieceList[idx] = Square(Empty)
	}

	for idx := range b.Material {
		b.Material[idx] = 0
	}

	for idx := range b.PieceNumber {
		b.PieceNumber[idx] = 0
	}

	b.SideToMove = Both
	b.EnPassant = NoSquare
	b.FiftyMove = 0
	b.HistoryPly = 0
	b.Ply = 0
	b.Castling = 0
	b.Hash = 0
	b.MoveList[b.Ply] = 0

}

func (b *Board) UpdatePieceList() {
	for idx := range 64 {
		sq := Fr64To120[idx]

		piece := b.Pieces[sq]

		if piece != Empty {
			color := pieceColor[piece]
			b.Material[color] += pieceValue[piece]
			b.PieceList[b.PieceIdx(piece)] = Square(sq)
			b.PieceNumber[piece]++
		}

	}
}

func (b *Board) String() string {
	var builder strings.Builder

	builder.WriteString("\nGame Board:\n\n")

	for rank := Rank8; rank >= Rank1; rank-- {
		fmt.Fprintf(&builder, "%d  ", rank+1)

		for file := FileA; file <= FileH; file++ {
			sq := FR2Square(file, rank)
			piece := b.Pieces[sq]

			fmt.Fprintf(&builder, "%c ", PieceChar[piece])
		}

		builder.WriteString("\n")
	}

	builder.WriteString("\n   a b c d e f g h\n\n")

	sideStr := ""

	switch b.SideToMove {
	case White:
		sideStr = "W"
	case Black:
		sideStr = "B"
	}

	fmt.Fprintf(&builder, "Side To Move: %s\n", sideStr)
	fmt.Fprintf(&builder, "Enpassant Square: %c%c\n", FileChar[FilesBoard[b.EnPassant]], RankChar[RanksBoard[b.EnPassant]])

	castleStr := ""

	if b.Castling&WKSide != 0 {
		castleStr += "K"
	}

	if b.Castling&WQSide != 0 {
		castleStr += "Q"
	}

	if b.Castling&BKSide != 0 {
		castleStr += "k"
	}

	if b.Castling&BQSide != 0 {
		castleStr += "q"
	}

	fmt.Fprintf(&builder, "Castling: %s\n", castleStr)
	fmt.Fprintf(&builder, "Hash: %x\n", b.Hash)

	return builder.String()
}
