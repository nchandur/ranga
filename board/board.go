package board

type Board struct {
	Pieces     []int8    // element represents the square on 8x8 board which is embedded in a 12x10 board
	Pawns      [3]uint64 // for white, black and both. bit will be set to 1 if a piece of that color exists on that square
	KingSq     [2]int    // square the kings are on (black and white)
	SideToMove int
	EnPass     int // target square of en passant
	FiftyMove  int // fifty move counter

	Ply     int
	HistPly int // how many plies have been made in the entire history of the game

	PositionKey uint64 // unique key generated for position

	CastlingPermission uint8   // stores integer that tells which side castling is permitted during the course of the game
	PieceNumber        [13]int // how many of each piece are on that board

	BigPieces   [3]int // any piece that isn't a pawn
	MajorPieces [3]int // rooks and queens
	MinorPieces [3]int // bishops and knights

	History []Undo // at whatever move number, what is possible
}

type Undo struct {
	Move               int
	CastlingPermission uint8
	EnPass             int
	FiftyMove          int
	PositionKey        uint64
}
