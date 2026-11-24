package board

// side-specific parameters
var (
	pawns      BitBoard
	promoMask  BitBoard // rank pawns land on when promoting
	doubleMask BitBoard // rank a single push lands on if a double push is possible
	push       int      // single push, in squares
	capL, capR int      // capture shifts toward the A and H files
)

// adds a quiet or capture move from source
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

// adds single pushes for every target bit
func (m *MoveList) addPawnPushes(targets BitBoard, delta int, promoMask BitBoard) {
	promos := targets & promoMask
	plain := targets &^ promoMask

	for promos != 0 {
		target := Square(promos.GetLSB())
		source := Square(int(target) - delta)
		m.AddMove(NewMove(source, target, MovePromoQ))
		m.AddMove(NewMove(source, target, MovePromoR))
		m.AddMove(NewMove(source, target, MovePromoB))
		m.AddMove(NewMove(source, target, MovePromoN))
		promos.PopBit(target)
	}

	for plain != 0 {
		target := Square(plain.GetLSB())
		m.AddMove(NewMove(Square(int(target)-delta), target, MoveQuiet))
		plain.PopBit(target)
	}
}

// adds captures for every target bit
func (m *MoveList) addPawnCaptures(targets BitBoard, delta int, promoMask BitBoard) {
	promos := targets & promoMask
	plain := targets &^ promoMask

	for promos != 0 {
		target := Square(promos.GetLSB())
		source := Square(int(target) - delta)
		m.AddMove(NewMove(source, target, MovePromoCaptureQ))
		m.AddMove(NewMove(source, target, MovePromoCaptureR))
		m.AddMove(NewMove(source, target, MovePromoCaptureB))
		m.AddMove(NewMove(source, target, MovePromoCaptureN))
		promos.PopBit(target)
	}

	for plain != 0 {
		target := Square(plain.GetLSB())
		m.AddMove(NewMove(Square(int(target)-delta), target, MoveCapture))
		plain.PopBit(target)
	}
}

// generates pawn moves
func (m *MoveList) generatePawnMoves(b *Board, capturesOnly bool) {

	if b.Side == White {
		pawns = b.PieceBitBoards[WP]
		promoMask, doubleMask = EigthRank, ThirdRank
		push, capL, capR = -8, -9, -7
	} else {
		pawns = b.PieceBitBoards[BP]
		promoMask, doubleMask = FirstRank, SixthRank
		push, capL, capR = 8, 7, 9
	}

	// quiet single/double pawn pushes
	if !capturesOnly {
		singles := shiftPawns(pawns, push) & ^b.Occupancies[Both]
		doubles := shiftPawns(singles&doubleMask, push) & ^b.Occupancies[Both]

		m.addPawnPushes(singles, push, promoMask)

		for doubles != 0 {
			target := Square(doubles.GetLSB())
			m.AddMove(NewMove(Square(int(target)-2*push), target, MoveDouble))
			doubles.PopBit(target)
		}
	}

	// standard pawn captures
	enemy := b.Occupancies[b.Side^1]
	m.addPawnCaptures(shiftPawns(pawns&NotAFile, capL)&enemy, capL, promoMask)
	m.addPawnCaptures(shiftPawns(pawns&NotHFile, capR)&enemy, capR, promoMask)

	// en-passant captures
	if b.EnPassant != NoSquare {
		ep := Square(b.EnPassant)
		attackers := PawnAttacks[b.Side^1][ep] & pawns
		for attackers != 0 {
			source := Square(attackers.GetLSB())
			m.AddMove(NewMove(source, ep, MoveEnPass))
			attackers.PopBit(source)
		}
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
