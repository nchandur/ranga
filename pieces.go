package main

type Piece int

const (
	Empty Piece = iota
	wP
	wN
	wB
	wR
	wQ
	wK
	bP
	bN
	bB
	bR
	bQ
	bK
)

// returns true if piece is pawn
var isPawn = []bool{false, true, false, false, false, false, false, true, false, false, false, false, false}

// returns true if piece is knight
var isKnight = []bool{false, false, true, false, false, false, false, false, true, false, false, false, false}

// returns true if piece is bishop
var isBishop = []bool{false, false, false, true, false, false, false, false, false, true, false, false, false}

// returns true if piece is rook
var isRook = []bool{false, false, false, false, true, false, false, false, false, false, true, false, false}

// returns true if piece is queen
var isQueen = []bool{false, false, false, false, false, true, false, false, false, false, false, true, false}

// returns true if piece is king
var isKing = []bool{false, false, false, false, false, false, true, false, false, false, false, false, true}

// returns true the piece slides across the board
var isSlider = []bool{false, false, false, true, true, true, false, false, false, true, true, true, false}

// returns true if piece is major
var isMajor = []bool{false, false, false, false, true, true, false, false, false, false, true, true, false}

// returns true if piece is minor
var isMinor = []bool{false, false, true, true, false, false, false, false, true, true, false, false, false}

// return material value of piece
var pieceValue = []int{0, 100, 300, 300, 500, 1000, 50000, 100, 300, 300, 500, 1000, 50000}

// returns color of piece
var pieceColor = []Color{Both, White, White, White, White, White, White, Black, Black, Black, Black, Black, Black}

// squares a knight can reach from it's position on a 120-index board
var knightDir = []Square{-21, -19, -12, -8, 8, 12, 19, 21}

// squares a bishop (or queen) can reach from it's position on a 120-index board
var bishopDir = []Square{-9, -11, 9, 11}

// squares a rook (or queen) can reach from it's position on a 120-index board
var rookDir = []Square{-10, -1, 1, 10}

// squares a king can reach from it's position on a 120-index board
var kingDir = []Square{-11, -10, -9, -1, 1, 9, 10, 11}

// number of directions each piece can move
var dirNum = []int{0, 0, 8, 4, 4, 8, 8, 0, 8, 4, 4, 8, 8}
var pieceDir = [][]Square{nil, nil, knightDir, bishopDir, rookDir, kingDir, kingDir, nil, knightDir, bishopDir, rookDir, kingDir, kingDir}

// non slider slices
var nonSlidePieceIdx = []int{0, 3}
var nonSlidePiece = []Piece{wN, wK, 0, bN, bK, 0}

// slider slices
var slidePieceIdx = []int{0, 4}
var slidePiece = []Piece{wB, wR, wQ, 0, bB, bR, bQ, 0}
