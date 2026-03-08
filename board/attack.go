package board

// check if given square is attacked by a pawn
func (b *Board) isAttackedByPawn(square Square, side Color) bool {
	switch side {
	case White:
		if (b.Pieces[square-11] == wP) || (b.Pieces[square-9] == wP) {
			return true
		}
	case Black:
		if (b.Pieces[square+11] == bP) || (b.Pieces[square+9] == bP) {
			return true
		}

	}
	return false
}

// check if given square is attacked by a knight
func (b *Board) isAttackedByKnight(square Square, side Color) bool {

	// squares a knight can reach from it's position on a 120-index board
	direction := []int{-8, -19, -21, -12, 8, 19, 21, 12}

	for _, dir := range direction {
		piece := b.Pieces[square+Square(dir)]

		if piece != Piece(Offboard) && piece.isKnight() && pieceColor[piece] == side {
			return true
		}

	}
	return false
}

// check if given square is attacked by rook or queen (horizonal and vertical)
func (b *Board) isAttackedByRookOrQueen(square Square, side Color) bool {

	// range of squares a rook can reach from it's position on a 120-index board
	direction := []int{-1, -10, 1, 10}

	for _, dir := range direction {
		tempSq := square + Square(dir)

		piece := b.Pieces[tempSq]

		// while square is not offboard
		for piece != 120 {

			// if there is a piece on square
			if piece != Empty {
				if (piece.isRook() || piece.isQueen()) && pieceColor[piece] == side {
					return true
				}
				break
			}
			tempSq += Square(dir)
			piece = b.Pieces[tempSq]

		}

	}

	return false
}

// check if given square is attacked by bishop or queen (diagonal)
func (b *Board) isAttackedByBishopOrQueen(square Square, side Color) bool {

	// range of squares a bishop can reach from it's position on a 120-based index
	direction := []int{-9, -11, 11, 9}

	for _, dir := range direction {
		tempSq := square + Square(dir)

		piece := b.Pieces[tempSq]

		// while square is not offboard
		for piece != 120 {

			// if there is a piece on square
			if piece != Empty {
				if (piece.isBishop() || piece.isQueen()) && pieceColor[piece] == side {
					return true
				}
				break
			}
			tempSq += Square(dir)
			piece = b.Pieces[tempSq]

		}

	}
	return false
}

// check if given square is attacked by a king
func (b *Board) isAttackedByKing(square Square, side Color) bool {

	// range of squares a king can reach from it's position on a 120-based index
	direction := []int{-1, -10, 1, 10, -9, -11, 11, 9}

	for _, dir := range direction {
		piece := b.Pieces[square+Square(dir)]

		if piece.isKing() && pieceColor[piece] == side {
			return true
		}

	}

	return false
}

func (b *Board) IsAttacked(square Square, side Color) bool {
	return b.isAttackedByPawn(square, side) || b.isAttackedByKnight(square, side) || b.isAttackedByBishopOrQueen(square, side) || b.isAttackedByRookOrQueen(square, side) || b.isAttackedByKing(square, side)
}
