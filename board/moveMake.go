package board

var castlePermission = []int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

// removes piece from board
func (b *Board) ClearPiece(piece Piece, square Square) {
	if square == 120 {
		return
	}

	pce := b.Pieces[square]

	if pce == Empty {
		return
	}

	color := pieceColor[piece]

	tempPieceNum := -1

	// hash out piece from position key
	b.hashPiece(piece, square)

	b.Pieces[square] = Empty
	b.Material[color] -= pieceValue[piece]

	if pieceBig[piece] {
		b.BigPieces[color]--

		if pieceMaj[piece] {
			b.MajorPieces[color]--
		}

		if pieceMin[piece] {
			b.MinorPieces[color]--
		}
	} else {
		b.Pawns[color].ClearBit(square)
		b.Pawns[Both].ClearBit(square)
	}

	for idx := range b.PieceNumber[pce] {
		if b.PieceList[pce][idx] == int(square) {
			tempPieceNum = idx
			break
		}
	}

	if tempPieceNum == -1 {
		return
	}

	b.PieceNumber[pce]--
	b.PieceList[pce][tempPieceNum] = b.PieceList[pce][b.PieceNumber[pce]]

}

func (b *Board) AddPiece(piece Piece, square Square) {

	if square == 120 {
		return
	}

	if piece == Empty {
		return
	}

	color := pieceColor[piece]

	b.hashPiece(piece, square)

	b.Pieces[square] = piece

	if pieceBig[piece] {
		b.BigPieces[color]++

		if pieceMaj[piece] {
			b.MajorPieces[color]++
		}

		if pieceMin[piece] {
			b.MinorPieces[color]++
		}
	} else {
		b.Pawns[color].SetBit(square)
		b.Pawns[Both].SetBit(square)
	}

	b.Material[color] += pieceValue[piece]
	b.PieceList[piece][b.PieceNumber[piece]] = int(square)
	b.PieceNumber[piece]++

}

func (b *Board) MovePiece(from, to Square) {

	if from == 120 || to == 120 {
		return
	}

	piece := b.Pieces[from]
	color := pieceColor[piece]

	b.hashPiece(piece, from)
	b.Pieces[from] = Empty

	b.hashPiece(piece, to)
	b.Pieces[to] = piece

	if !pieceBig[piece] {
		b.Pawns[color].ClearBit(from)
		b.Pawns[Both].ClearBit(from)
		b.Pawns[color].SetBit(to)
		b.Pawns[Both].SetBit(to)
	}

	for idx := range b.PieceNumber[piece] {
		if b.PieceList[piece][idx] == int(from) {
			b.PieceList[piece][idx] = int(to)
			break
		}
	}

}
