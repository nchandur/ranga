package board

// side-specific parameters
var (
	pawns            BitBoard
	push             int
	startLo, startHi Square
	promoLo, promoHi Square
)

// adds quiet or capture move from source
func (m *MoveList) addTargets(b *Board, source Square, attacks BitBoard) {
	for attacks != 0 {
		target := Square(attacks.GetLSB())

		if b.Occupancies[b.Side^1].GetBit(target) == 0 {
			m.AddMove(NewMove(source, target, MoveQuiet))
		} else {
			m.AddMove(NewMove(source, target, MoveCapture))
		}

		attacks.PopBit(target)
	}
}

// generates pseudo-legal moves for every piece in pieces
func (m *MoveList) addPieceMoves(b *Board, pieces BitBoard, capturesOnly bool,
	attackFn func(Square, BitBoard) BitBoard) {

	for pieces != 0 {
		source := Square(pieces.GetLSB())

		// excludes friendly blockers
		attacks := attackFn(source, b.Occupancies[Both]) & ^b.Occupancies[b.Side]
		if capturesOnly {
			attacks &= b.Occupancies[b.Side^1]
		}

		m.addTargets(b, source, attacks)
		pieces.PopBit(source)
	}
}

func knightAttacks(s Square, _ BitBoard) BitBoard {
	return KnightAttacks[s]
}

// generates pawn moves
func (m *MoveList) generatePawnMoves(b *Board, capturesOnly bool) {

	if b.Side == White {
		pawns = b.PieceBitBoards[WP]
		push = -8
		startLo, startHi = A2, H2
		promoLo, promoHi = A7, H7
	} else {
		pawns = b.PieceBitBoards[BP]
		push = 8
		startLo, startHi = A7, H7
		promoLo, promoHi = A2, H2
	}

	occ := b.Occupancies[Both]

	// en-passant target as a bitboard
	var epBB BitBoard
	if b.EnPassant != NoSquare {
		epBB = BitBoard(uint64(1) << uint64(b.EnPassant))
	}

	for pawns != 0 {
		source := Square(pawns.GetLSB())
		target := Square(int(source) + push)
		promo := source >= promoLo && source <= promoHi

		// quiet single/double pawn push
		if !capturesOnly && target >= A8 && target <= H1 && occ.GetBit(target) == 0 {
			if promo {
				m.AddMove(NewMove(source, target, MovePromoQ))
				m.AddMove(NewMove(source, target, MovePromoR))
				m.AddMove(NewMove(source, target, MovePromoB))
				m.AddMove(NewMove(source, target, MovePromoN))
			} else {
				m.AddMove(NewMove(source, target, MoveQuiet))
				// double push from starting rank
				if source >= startLo && source <= startHi {
					doubleTarget := Square(int(target) + push)
					if occ.GetBit(doubleTarget) == 0 {
						m.AddMove(NewMove(source, doubleTarget, MoveDouble))
					}
				}
			}
		}

		// standard pawn captures
		attacks := PawnAttacks[b.Side][source] & b.Occupancies[b.Side^1]

		for attacks != 0 {
			target = Square(attacks.GetLSB())
			if promo {
				m.AddMove(NewMove(source, target, MovePromoCaptureQ))
				m.AddMove(NewMove(source, target, MovePromoCaptureR))
				m.AddMove(NewMove(source, target, MovePromoCaptureB))
				m.AddMove(NewMove(source, target, MovePromoCaptureN))
			} else {
				m.AddMove(NewMove(source, target, MoveCapture))
			}
			attacks.PopBit(target)
		}

		// en-passant captures
		if epBB != 0 {
			if epAttacks := PawnAttacks[b.Side][source] & epBB; epAttacks != 0 {
				m.AddMove(NewMove(source, Square(epAttacks.GetLSB()), MoveEnPass))
			}
		}

		pawns.PopBit(source)
	}
}

// generate knight moves
func (m *MoveList) generateKnightMoves(b *Board, capturesOnly bool) {
	pce := BN
	if b.Side == White {
		pce = WN
	}
	m.addPieceMoves(b, b.PieceBitBoards[pce], capturesOnly, knightAttacks)
}

// generate king moves
func (m *MoveList) generateKingMoves(b *Board, capturesOnly bool) {
	occ := b.Occupancies[Both]

	pce := BK
	if b.Side == White {
		pce = WK
	}

	if king := b.PieceBitBoards[pce]; king != 0 {
		source := Square(king.GetLSB())

		// excludes friendly blockers
		attacks := KingAttacks[source] & ^b.Occupancies[b.Side]
		if capturesOnly {
			attacks &= b.Occupancies[b.Side^1]
		}
		m.addTargets(b, source, attacks)
	}

	if capturesOnly {
		return
	}

	if b.Side == White {
		// king-side
		if b.Castle&WKCA != 0 && occ.GetBit(F1) == 0 && occ.GetBit(G1) == 0 {
			if !b.IsSquareAttacked(E1, Black) && !b.IsSquareAttacked(F1, Black) {
				m.AddMove(NewMove(E1, G1, MoveCastle))
			}
		}
		// queen-side
		if b.Castle&WQCA != 0 && occ.GetBit(D1) == 0 && occ.GetBit(C1) == 0 && occ.GetBit(B1) == 0 {
			if !b.IsSquareAttacked(E1, Black) && !b.IsSquareAttacked(D1, Black) {
				m.AddMove(NewMove(E1, C1, MoveCastle))
			}
		}
	} else {
		// king-side
		if b.Castle&BKCA != 0 && occ.GetBit(F8) == 0 && occ.GetBit(G8) == 0 {
			if !b.IsSquareAttacked(E8, White) && !b.IsSquareAttacked(F8, White) {
				m.AddMove(NewMove(E8, G8, MoveCastle))
			}
		}
		// queen-side
		if b.Castle&BQCA != 0 && occ.GetBit(D8) == 0 && occ.GetBit(C8) == 0 && occ.GetBit(B8) == 0 {
			if !b.IsSquareAttacked(E8, White) && !b.IsSquareAttacked(D8, White) {
				m.AddMove(NewMove(E8, C8, MoveCastle))
			}
		}
	}
}

// generates bishop moves
func (m *MoveList) generateBishopMoves(b *Board, capturesOnly bool) {
	pce := BB
	if b.Side == White {
		pce = WB
	}
	m.addPieceMoves(b, b.PieceBitBoards[pce], capturesOnly, GetBishopAttacks)
}

// generates pseudo-legal rook moves
func (m *MoveList) generateRookMoves(b *Board, capturesOnly bool) {
	pce := BR
	if b.Side == White {
		pce = WR
	}
	m.addPieceMoves(b, b.PieceBitBoards[pce], capturesOnly, GetRookAttacks)
}

// generates pseudo-legal queen moves
func (m *MoveList) generateQueenMoves(b *Board, capturesOnly bool) {
	pce := BQ
	if b.Side == White {
		pce = WQ
	}
	m.addPieceMoves(b, b.PieceBitBoards[pce], capturesOnly, GetQueenAttacks)
}

// generate pseudo-legal moves for all pieces
func (m *MoveList) generate(b *Board, capturesOnly bool) {
	m.generatePawnMoves(b, capturesOnly)
	m.generateKnightMoves(b, capturesOnly)
	m.generateKingMoves(b, capturesOnly)
	m.generateBishopMoves(b, capturesOnly)
	m.generateRookMoves(b, capturesOnly)
	m.generateQueenMoves(b, capturesOnly)
}

// generate all pseudo-legal moves
func (m *MoveList) GenerateMoves(b *Board) {
	m.generate(b, false)
}

// generate only capture moves
func (m *MoveList) GenerateCaptures(b *Board) {
	m.generate(b, true)
}
