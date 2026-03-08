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
func (b *Board) ClearPiece(square Square) {
	if square == 120 {
		return
	}

	pce := b.Pieces[square]

	if pce == Empty {
		return
	}

	color := pieceColor[pce]

	tempPieceNum := -1

	// hash out piece from position key
	b.hashPiece(pce, square)

	b.Pieces[square] = Empty
	b.Material[color] -= pieceValue[pce]

	if pieceBig[pce] {
		b.BigPieces[color]--

		if pieceMaj[pce] {
			b.MajorPieces[color]--
		}

		if pieceMin[pce] {
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

func (b *Board) MakeMove(move Move) bool {

	if !b.Check() {
		return false
	}

	from := Square(move.FromSq())
	to := Square(move.ToSq())

	side := b.SideToMove

	if from == NoSquare || from == Offboard || to == NoSquare || to == Offboard {
		return false
	}

	if side != White && side != Black {
		return false
	}

	if b.Pieces[from] == Empty {
		return false
	}

	b.History[b.HistPly].PositionKey = b.PositionKey

	if move.IsEnPassant() {
		switch side {
		case White:
			b.ClearPiece(to - 10)
		case Black:
			b.ClearPiece(to + 10)
		}

	}

	if move.IsCastle() {
		switch to {
		case C1:
			b.MovePiece(A1, D1)
		case C8:
			b.MovePiece(A8, D8)
		case G1:
			b.MovePiece(H1, F1)
		case G8:
			b.MovePiece(H8, F8)
		default:
			return false
		}
	}

	if b.EnPass != NoSquare {
		b.hashEnPassant()
	}

	b.hashCastle()
	b.History[b.HistPly].Move = move.move
	b.History[b.HistPly].FiftyMove = b.FiftyMove
	b.History[b.HistPly].EnPass = b.EnPass
	b.History[b.HistPly].CastlingPermission = b.CastlingPermission

	b.CastlingPermission &= Castling(castlePermission[from])
	b.CastlingPermission &= Castling(castlePermission[to])
	b.EnPass = NoSquare
	b.hashCastle()

	captured := move.CapturedPiece()
	b.FiftyMove++

	if captured != int(Empty) {
		b.ClearPiece(to)
		b.FiftyMove = 0
	}

	b.HistPly++
	b.Ply++

	if b.Pieces[from].isPawn() {
		b.FiftyMove = 0

		if move.IsPawnStart() {
			switch side {
			case White:
				b.EnPass = from + 10
				if Fr120ToRank(int(b.EnPass)) != Three {
					return false
				}
			case Black:
				b.EnPass = from - 10
				if Fr120ToRank(int(b.EnPass)) != Six {
					return false
				}
			}
			b.hashEnPassant()
		}

	}

	b.MovePiece(from, to)

	promoted := Piece(move.PromotedPiece())

	if promoted != Empty {
		if promoted.isPawn() {
			return false
		}

		b.ClearPiece(to)

		b.AddPiece(promoted, to)

	}

	if b.Pieces[to].isKing() {
		b.KingSq[b.SideToMove] = to
	}

	b.SideToMove ^= 1

	b.hashSideToMove()

	if b.IsAttacked(b.KingSq[side], b.SideToMove) {
		// TakeMove
		return false
	}

	return b.Check()
}
