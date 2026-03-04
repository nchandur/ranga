package board

import "fmt"

type Board struct {
	Pieces     []Piece  // element represents the square on 8x8 board which is embedded in a 12x10 board
	Pawns      []uint64 // for white, black and both. bit will be set to 1 if a piece of that color exists on that square
	KingSq     []Square // square the kings are on (black and white)
	SideToMove Color
	EnPass     Square // target square of en passant
	FiftyMove  int    // fifty move counter

	Ply     int
	HistPly int // how many plies have been made in the entire history of the game

	PositionKey uint64 // unique key generated for position

	CastlingPermission Castling // stores integer that tells which side castling is permitted during the course of the game
	PieceNumber        []int    // how many of each piece are on that board

	BigPieces   []int // any piece that isn't a pawn
	MajorPieces []int // rooks and queens
	MinorPieces []int // bishops and knights
	Material    []int

	History []Undo // at whatever move number, what is possible

	PieceList [][]int // for each piece type there can be up to 10 of it on the board

}

type Undo struct {
	Move               int
	CastlingPermission uint8
	EnPass             int
	FiftyMove          int
	PositionKey        uint64
}

func NewBoard() Board {
	b := Board{}

	b.Pieces = make([]Piece, 120)
	b.Pawns = make([]uint64, 3)
	b.KingSq = make([]Square, 2)
	b.SideToMove = Both
	b.EnPass = NoSquare
	b.FiftyMove = 0
	b.Ply = 0
	b.HistPly = 0
	b.PositionKey = uint64(0)
	b.CastlingPermission = Castling(0)
	b.PieceNumber = make([]int, 13)

	b.BigPieces = make([]int, 2)
	b.MajorPieces = make([]int, 2)
	b.MinorPieces = make([]int, 2)
	b.Material = make([]int, 2)

	b.PieceList = make([][]int, 13)

	for i := range 13 {
		b.PieceList[i] = make([]int, 10)
	}

	return b
}

// update piece list on board
func (b *Board) UpdatePieceList() {

	// all non-pawn and non-empty pieces ".PNBRQKpnbrqk"
	pieceBig := []bool{false, false, true, true, true, true, true, false, true, true, true, true, true}

	// all major pieces (rooks and queens and kings)
	pieceMaj := []bool{false, false, false, false, true, true, true, false, false, false, true, true, true}

	// all minor pieces (knights and bishops)
	pieceMin := []bool{false, false, true, true, false, false, false, false, true, true, false, false, false}

	// numerical value of each piece
	pieceValue := []int{0, 100, 300, 300, 500, 1000, 50000, 100, 300, 300, 500, 1000, 50000}

	// color of each piece
	pieceColor := []Color{Both, White, White, White, White, White, White, Black, Black, Black, Black, Black, Black}

	for idx := range 120 {
		piece := b.Pieces[idx]

		if piece != 120 && piece != Empty {
			color := pieceColor[piece]

			// increment the big pieces
			if pieceBig[piece] {
				b.BigPieces[color]++
			}

			// increment major pieces
			if pieceMaj[piece] {
				b.MajorPieces[color]++
			}

			// increment minor pieces
			if pieceMin[piece] {
				b.MinorPieces[color]++
			}

			b.Material[color] += pieceValue[piece]

			b.PieceList[piece][b.PieceNumber[piece]] = idx
			b.PieceNumber[piece]++

		}

	}

}

// generates unique position key for the board
func (b *Board) GenPositionKey() uint64 {

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
	if b.EnPass != NoSquare {
		key ^= EnPassKeys[b.EnPass]
	}

	// castling rights
	key ^= CastleKeys[b.CastlingPermission]

	return key

}

func (b *Board) Reset() {

	// make all squares off-board
	for i := range 120 {
		b.Pieces[i] = 120
	}

	// reset all 8x8 squares to empty
	for i := range 64 {
		b.Pieces[Fr64To120(i)] = Empty
	}

	for i := range 3 {
		b.BigPieces[i] = 0
		b.MajorPieces[i] = 0
		b.MinorPieces[i] = 0
		b.Pawns[i] = uint64(0)
	}

	for i := range 13 {
		b.PieceNumber[i] = 0
	}

	b.KingSq[0] = NoSquare
	b.KingSq[1] = NoSquare

	b.SideToMove = Both
	b.EnPass = NoSquare
	b.FiftyMove = 0

	b.Ply = 0
	b.HistPly = 0
	b.CastlingPermission = 0
	b.PositionKey = uint64(0)

	b.History = nil

}

// print board
func (b *Board) Print() {

	pieces := ".PNBRQKpnbrqk"
	side := "wb-"

	fmt.Printf("\n===========Board===========\n\n")

	for rank := int(Eight); rank >= int(One); rank-- {
		fmt.Printf("%d ", rank+1)
		for file := A; file <= H; file++ {
			sq := FRToSq(File(file), Rank(rank))
			piece := b.Pieces[sq]
			fmt.Printf("%3c", pieces[piece])

		}
		fmt.Println()
	}

	fmt.Printf("\n  ")
	for file := A; file <= H; file++ {
		fmt.Printf("%3c", 'a'+file)
	}

	fmt.Printf("\n===========================\n\n")

	fmt.Println()
	fmt.Printf("Side: %c\n", side[b.SideToMove])
	fmt.Printf("En Passant Square: %s\n", Fr120ToFR(int(b.EnPass)))
	fmt.Printf("Castle: %x%x%x%x\n", b.CastlingPermission&WKSide, b.CastlingPermission&WQSide, b.CastlingPermission&BKSide, b.CastlingPermission&BQSide)
	fmt.Printf("Position Key: %x\n", b.PositionKey)

}
