package main

// check if given square is attacked by a pawn
func (b *Board) isAttackedByPawn(square Square, side Color) bool {
	switch side {
	case White:
		if b.Pieces[square-11] == wP || b.Pieces[square-9] == wP {
			return true
		}
	case Black:
		if b.Pieces[square+11] == bP || b.Pieces[square+9] == bP {
			return true
		}
	}
	return false
}

// check if given square is attacked by a knight
func (b *Board) isAttackedByKnight(square Square, side Color) bool {

	if b.Pieces[square] == Piece(Offboard) {
		return false
	}

	if pieceColor[b.Pieces[square]] == side {
		return false
	}

	for idx := range knightDir {
		pce := b.Pieces[square+knightDir[idx]]

		if pce == Piece(Offboard) {
			continue
		}

		if b.PieceList[pce] != Offboard && pieceColor[pce] == side && isKnight[pce] {
			return true
		}

	}

	return false
}

// check if given square is attacked by rook or queen (horizontal and vertical)
func (b *Board) isAttackedByRookOrQueen(square Square, side Color) bool {

	if b.Pieces[square] == Piece(Offboard) {
		return false
	}

	if pieceColor[b.Pieces[square]] == side {
		return false
	}

	for _, dir := range rookDir {
		tSq := square + dir

		pce := b.Pieces[tSq]

		for pce != Piece(Offboard) {
			if pce != Empty {
				if (isRook[pce] || isQueen[pce]) && pieceColor[pce] == side {
					return true
				}
				break
			}
			tSq += dir
			pce = b.Pieces[tSq]
		}

	}

	return false
}

// check if given square is attacked by bishop or queen (diagonal)
func (b *Board) isAttackedByBishopOrQueen(square Square, side Color) bool {

	if b.Pieces[square] == Piece(Offboard) {
		return false
	}

	if pieceColor[b.Pieces[square]] == side {
		return false
	}

	for _, dir := range bishopDir {
		tSq := square + dir

		pce := b.Pieces[tSq]

		for pce != Piece(Offboard) {
			if pce != Empty {
				if (isBishop[pce] || isQueen[pce]) && pieceColor[pce] == side {
					return true
				}
				break
			}
			tSq += dir
			pce = b.Pieces[tSq]
		}

	}

	return false

}

// check if given square is attacked by a king
func (b *Board) isAttackedByKing(square Square, side Color) bool {

	if b.Pieces[square] == Piece(Offboard) {
		return false
	}

	if pieceColor[b.Pieces[square]] == side {
		return false
	}

	for idx := range kingDir {
		pce := b.Pieces[square+kingDir[idx]]

		if pce == Piece(Offboard) {
			continue
		}

		if b.PieceList[pce] != Offboard && pieceColor[pce] == side && isKing[pce] {
			return true
		}

	}

	return false
}

func (b *Board) IsAttacked(square Square, side Color) bool {
	return b.isAttackedByPawn(square, side) ||
		b.isAttackedByKnight(square, side) ||
		b.isAttackedByBishopOrQueen(square, side) ||
		b.isAttackedByRookOrQueen(square, side) ||
		b.isAttackedByKing(square, side)
}
