package main

var kings = []Piece{wK, bK}

var castlePermissions = []int{
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


func (b *Board) MakeMove(move Move) bool {
	from := move.FromSquare()
	to := move.ToSquare()
	side := b.SideToMove

	b.History[b.HistoryPly].Hash = b.Hash

	if move&MFLAGEP != 0 {
		if side == White {
			b.clearPiece(to - 10)
		} else {
			b.clearPiece(to + 10)
		}
	} else if move&MFLAGCA != 0 {
		switch to {
		case C1:
			b.movePiece(A1, D1)
		case C8:
			b.movePiece(A8, D8)
		case G1:
			b.movePiece(H1, F1)
		case G8:
			b.movePiece(H8, F8)
		}
	}

	if b.EnPassant != NoSquare {
		b.hashEnpassant()
	}

	b.hashCastle()

	b.History[b.HistoryPly].Move = move
	b.History[b.HistoryPly].FiftyMove = b.FiftyMove
	b.History[b.HistoryPly].EnPassant = b.EnPassant
	b.History[b.HistoryPly].CastleBit = b.CastleBit

	b.CastleBit &= CastleBit(castlePermissions[from])
	b.CastleBit &= CastleBit(castlePermissions[to])
	b.EnPassant = NoSquare

	b.hashCastle()

	captured := move.Captured()
	b.FiftyMove++

	if captured != Empty {
		b.clearPiece(to)
		b.FiftyMove = 0
	}

	b.HistoryPly++
	b.Ply++

	if isPawn[b.Pieces[from]] {
		b.FiftyMove = 0
		if move&MFLAGPS != 0 {
			if side == White {
				b.EnPassant = from + 10
			} else {
				b.EnPassant = from - 10
			}
			b.hashEnpassant()
		}
	}

	b.movePiece(from, to)

	promoted := move.Promoted()

	if promoted != Empty {
		b.clearPiece(to)
		b.addPiece(promoted, to)
	}

	b.SideToMove ^= 1
	b.hashSideToMove()

	if b.IsAttacked(b.PieceList[(int(kings[side])*10)+0], b.SideToMove) {
		b.TakeMove()
		return b.Checkboard()
	}

	return b.Checkboard()
}

func (b *Board) TakeMove() {
	b.HistoryPly--
	b.Ply--

	move := b.History[b.HistoryPly].Move
	from := move.FromSquare()
	to := move.ToSquare()

	if b.EnPassant != NoSquare {
		b.hashEnpassant()
	}

	b.hashCastle()
	b.CastleBit = b.History[b.HistoryPly].CastleBit
	b.FiftyMove = b.History[b.HistoryPly].FiftyMove
	b.EnPassant = b.History[b.HistoryPly].EnPassant

	if b.EnPassant != NoSquare {
		b.hashEnpassant()
	}
	b.hashCastle()

	b.SideToMove ^= 1
	b.hashSideToMove()

	if move&MFLAGEP != 0 {
		if b.SideToMove == White {
			b.addPiece(bP, to+10)
		} else {
			b.addPiece(wP, to-10)
		}
	} else if move&MFLAGCA != 0 {
		switch to {
		case C1:
			b.movePiece(D1, A1)
		case C8:
			b.movePiece(D8, A8)
		case G1:
			b.movePiece(F1, H1)
		case G8:
			b.movePiece(F8, H8)
		}
	}

	b.movePiece(to, from)

	captured := move.Captured()

	if captured != Empty {
		b.addPiece(captured, to)
	}

	if move.Promoted() != Empty {
		b.clearPiece(from)

		if pieceColor[move.Promoted()] == White {
			b.addPiece(wP, from)
		} else {
			b.addPiece(bP, from)
		}

	}

}
