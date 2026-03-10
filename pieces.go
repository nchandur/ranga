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
