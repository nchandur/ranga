package board

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

	b.BigPieces = make([]int, 3)
	b.MajorPieces = make([]int, 3)
	b.MinorPieces = make([]int, 3)

	b.PieceList = make([][]int, 13)

	for i := range 13 {
		b.PieceList[i] = make([]int, 10)
	}

	return b
}

func (b *Board) Reset() {

	// make all squares off-board
	for i := range 120 {
		b.Pieces[i] = 120
	}

	// reset all 8x8 squares to empty
	for i := range 64 {
		b.Pieces[Fr64To120(uint8(i))] = Empty
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
