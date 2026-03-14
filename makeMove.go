package main

func (b *Board) clearPiece(square Square) {
	piece := b.Pieces[square]
	col := pieceColor[piece]

	tPceNum := -1

	b.hashPiece(piece, square)

	b.Pieces[square] = Empty
	b.Material[col] -= pieceValue[piece]

	for idx := range b.PieceNumber[piece] {
		if b.PieceList[(int(piece)*10)+idx] == square {
			tPceNum = idx
			break
		}

	}

	b.PieceNumber[piece]--
	b.PieceList[(int(piece)*10)+tPceNum] = b.PieceList[(int(piece)*10)+b.PieceNumber[piece]]

}

func (b *Board) addPiece(piece Piece, square Square) {
	col := pieceColor[piece]

	b.hashPiece(piece, square)

	b.Pieces[square] = piece
	b.Material[col] += pieceValue[piece]
	b.PieceList[(int(piece)*10)+b.PieceNumber[piece]] = square
	b.PieceNumber[piece]++

}

func (b *Board) movePiece(from, to Square) {
	piece := b.Pieces[from]

	b.hashPiece(piece, from)
	b.Pieces[from] = Empty

	b.hashPiece(piece, to)
	b.Pieces[to] = piece

	for idx := range b.PieceNumber[piece] {

		if b.PieceList[(int(piece)*10)+idx] == from {
			b.PieceList[(int(piece)*10)+idx] = to
			break
		}
	}
}
