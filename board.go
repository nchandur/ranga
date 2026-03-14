package main

import (
	"fmt"
	"log"
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
	Pieces     []Piece // represents pieces on the 120-index board. the square on which a piece exists has its corresponding value
	SideToMove Color
	FiftyMove  int
	Ply        int
	HistoryPly int
	History    []struct {
		Move
		Hash      uint64
		EnPassant Square
		FiftyMove int
		CastleBit
	}
	CastleBit
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
	board.Material = make([]int, 2)

	board.MoveList = make([]int, MAX_DEPTH*MAX_POSITION_MOVES)
	board.MoveScores = make([]int, MAX_DEPTH*MAX_POSITION_MOVES)
	board.MoveListStart = make([]int, MAX_DEPTH)

	for range MAX_GAME_MOVES {
		board.History = append(board.History, struct {
			Move
			Hash      uint64
			EnPassant Square
			FiftyMove int
			CastleBit
		}{
			Move:      0,
			CastleBit: 0,
			EnPassant: NoSquare,
			FiftyMove: 0,
			Hash:      0,
		})
	}

	board.Reset()

	return board
}

func (b *Board) PieceIdx(piece Piece) int {
	return (int(piece) * 10) + b.PieceNumber[int(piece)]
}

func (b *Board) hashPiece(piece Piece, square Square) {
	b.Hash ^= PieceKeys[(int(piece)*120)+int(square)]
}

func (b *Board) hashCastle() {
	b.Hash ^= CastleKeys[b.CastleBit]
}

func (b *Board) hashEnpassant() {
	b.Hash ^= PieceKeys[b.EnPassant]
}

func (b *Board) hashSideToMove() {
	b.Hash ^= SideKey
}

func (b *Board) GenHash() uint64 {
	var res uint64

	for sq := range BOARD_SQUARE_NUM {
		piece := b.Pieces[sq]

		if piece != Empty && piece != Piece(Offboard) {
			res ^= PieceKeys[(int(piece)*120)+sq]
		}
	}

	if b.SideToMove == White {
		res ^= SideKey
	}

	if b.EnPassant != NoSquare {
		res ^= PieceKeys[b.EnPassant]
	}

	res ^= CastleKeys[b.CastleBit]

	return res
}

func (b *Board) Reset() {

	for idx := range BOARD_SQUARE_NUM {
		b.Pieces[idx] = 120
	}

	for idx := range 64 {
		b.Pieces[Fr64To120[idx]] = Empty
	}

	for idx := range b.PieceList {
		b.PieceList[idx] = NoSquare
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
	b.CastleBit = 0
	b.Hash = 0
	b.MoveList[b.Ply] = 0

}

func (b *Board) Checkboard() bool {

	tPieceNum := make([]int, 13)
	tMaterial := make([]int, 2)

	for tPce := wP; tPce <= bK; tPce++ {
		for tPceNum := range b.PieceNumber[tPce] {
			sq120 := b.PieceList[(int(tPce)*10)+tPceNum]
			if b.Pieces[sq120] != tPce {
				log.Println("piece mismatch")
				return false
			}
		}
	}

	for sq64 := range 64 {
		sq120 := Fr64To120[sq64]
		tPce := b.Pieces[sq120]
		tPieceNum[tPce]++

		if tPce != Empty {
			tMaterial[pieceColor[tPce]] += pieceValue[tPce]
		}
	}

	for tPce := wP; tPce <= bK; tPce++ {
		if tPieceNum[tPce] != b.PieceNumber[tPce] {
			log.Println("piece number mismatch")
			return false
		}
	}

	if tMaterial[White] != b.Material[White] || tMaterial[Black] != b.Material[Black] {
		log.Println("material mismatch")
		return false
	}

	if b.SideToMove != White && b.SideToMove != Black {
		log.Println("side to move mismatch")
		return false
	}

	if b.GenHash() != b.Hash {
		log.Println("hash mismatch")
		return false
	}

	return true
}

func (b *Board) UpdatePieceList() {

	for i := range b.Material {
		b.Material[i] = 0
	}

	for i := range b.PieceNumber {
		b.PieceNumber[i] = 0
	}

	for idx := range 64 {

		sq := Fr64To120[idx]
		piece := b.Pieces[sq]

		if piece != Empty {
			color := pieceColor[piece]
			b.Material[color] += pieceValue[piece]
			b.PieceList[(int(piece)*10)+b.PieceNumber[int(piece)]] = Square(sq)
			b.PieceNumber[int(piece)]++
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
		sideStr = "w"
	case Black:
		sideStr = "b"
	}

	fmt.Fprintf(&builder, "To Move: %s\n", sideStr)

	if b.EnPassant != NoSquare {
		fmt.Fprintf(&builder, "Enpassant Square: %c%c\n", FileChar[FilesBoard[b.EnPassant]], RankChar[RanksBoard[b.EnPassant]])
	} else {
		fmt.Fprintf(&builder, "Enpassant Square: -\n")
	}

	castleStr := ""

	if b.CastleBit&WKSide != 0 {
		castleStr += "K"
	}

	if b.CastleBit&WQSide != 0 {
		castleStr += "Q"
	}

	if b.CastleBit&BKSide != 0 {
		castleStr += "k"
	}

	if b.CastleBit&BQSide != 0 {
		castleStr += "q"
	}

	if castleStr == "" {
		fmt.Fprintf(&builder, "Castling: -\n")
	} else {
		fmt.Fprintf(&builder, "Castling: %s\n", castleStr)
	}
	fmt.Fprintf(&builder, "Hash: %x\n", b.Hash)

	return builder.String()
}
