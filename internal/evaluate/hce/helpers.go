package hce

import "ranga/internal/board"

func (h *HCE) evalPawns(b *board.Board, square board.Square, side board.Color) (int, int) {

	// calculates eval score for pawns
	open, end := 0, 0

	sq := square
	pce := board.WP
	passedPawns := WhitePassedPawnMasks[square] & b.PieceBitBoards[board.BP]

	if side == board.Black {
		sq = square ^ 56
		pce = board.BP
		passedPawns = BlackPassedPawnMasks[square] & b.PieceBitBoards[board.WP]
	}

	// increment piece square table evaluation
	open += PositionalScores[Opening][0][sq]
	end += PositionalScores[EndGame][0][sq]

	// double pawn
	if (b.PieceBitBoards[pce] & FileMasks[sq]).CountBits() > 1 {
		open += doublePawnPenalty[Opening]
		end += doublePawnPenalty[EndGame]
	}

	// isolated pawn
	if (b.PieceBitBoards[pce] & IsolatedMasks[sq]) == 0 {
		open += isolatedPawnPenalty[Opening]
		end += isolatedPawnPenalty[EndGame]
	}

	// passed pawn
	if passedPawns == 0 {
		open += passedPawnBonus[Opening][Ranks[sq]]
		end += passedPawnBonus[EndGame][Ranks[sq]]
	}

	return open, end
}

func (h *HCE) evalKnights(square board.Square, side board.Color) (int, int) {
	open, end := 0, 0

	sq := square

	if side == board.Black {
		sq = square ^ 56
	}

	// increment piece square table evaluation
	open += PositionalScores[Opening][1][sq]
	end += PositionalScores[EndGame][1][sq]

	return open, end
}

func (h *HCE) evalBishops(b *board.Board, square board.Square, side board.Color) (int, int) {
	open, end := 0, 0

	sq := square

	if side == board.Black {
		sq = square ^ 56
	}

	// increment piece square table evaluation
	open += PositionalScores[Opening][2][sq]
	end += PositionalScores[EndGame][2][sq]

	// increment mobility bonus
	open += board.GetBishopAttacks(square, b.Occupancies[board.Both]).CountBits() * bishopMobility[Opening]
	end += board.GetBishopAttacks(square, b.Occupancies[board.Both]).CountBits() * bishopMobility[EndGame]

	return open, end
}

func (h *HCE) evalRooks(b *board.Board, square board.Square, side board.Color) (int, int) {
	open, end := 0, 0

	sq := square
	semiopenfile := (b.PieceBitBoards[board.WP] & FileMasks[square])

	if side == board.Black {
		sq = square ^ 56
		semiopenfile = (b.PieceBitBoards[board.BP] & FileMasks[square])
	}

	// increment piece square table evaluation
	open += PositionalScores[Opening][3][sq]
	end += PositionalScores[EndGame][3][sq]

	// semi open file bonus
	if semiopenfile == 0 {
		open += semiOpenFile[Opening]
		end += semiOpenFile[EndGame]
	}
	// open file bonus
	if ((b.PieceBitBoards[board.WP] | b.PieceBitBoards[board.BP]) & FileMasks[square]) == 0 {
		open += openFile[Opening]
		end += openFile[EndGame]
	}

	return open, end
}

func (h *HCE) evalQueens(b *board.Board, square board.Square, side board.Color) (int, int) {
	open, end := 0, 0

	sq := square

	if side == board.Black {
		sq = square ^ 56
	}

	// increment piece square table evaluation
	open += PositionalScores[Opening][4][sq]
	end += PositionalScores[EndGame][4][sq]

	// increment mobility bonus
	open += board.GetQueenAttacks(square, b.Occupancies[board.Both]).CountBits() * queenMobility[Opening]
	end += board.GetQueenAttacks(square, b.Occupancies[board.Both]).CountBits() * queenMobility[EndGame]

	return open, end
}

func (h *HCE) evalKings(b *board.Board, square board.Square, side board.Color) (int, int) {
	open, end := 0, 0

	sq := square
	filesafetypenalties := (b.PieceBitBoards[board.WP] & FileMasks[square])
	shieldbonus := (board.KingAttacks[square] & b.Occupancies[board.White])

	if side == board.Black {
		sq = square ^ 56
		filesafetypenalties = (b.PieceBitBoards[board.BP] & FileMasks[square])
		shieldbonus = (board.KingAttacks[square] & b.Occupancies[board.Black])
	}

	// increment piece square table evaluation
	open += PositionalScores[Opening][5][sq]
	end += PositionalScores[EndGame][5][sq]

	// file safety penalties
	if filesafetypenalties == 0 {
		open -= semiOpenFile[Opening]
		end -= semiOpenFile[EndGame]
	}

	if ((b.PieceBitBoards[board.WP] | b.PieceBitBoards[board.BP]) & FileMasks[square]) == 0 {
		open -= openFile[Opening]
		end -= openFile[EndGame]
	}

	// shield bonus
	open += (shieldbonus.CountBits() * kingShieldBonus[Opening])
	end += (shieldbonus.CountBits() * kingShieldBonus[EndGame])

	return open, end
}

// calculates tapered evaluation
func (h *HCE) taper(open, end, phase, phaseScore int) int {

	score := 0

	// taper score between opening and endgame phases
	switch phase {
	case MiddleGame:
		score = ((open * phaseScore) + (end * (PhaseScore[Opening] - phaseScore))) / PhaseScore[Opening]
	case Opening:
		score = open
	case EndGame:
		score = end
	}

	return score
}

// calculates game phase value based on non-pawn material remaining on the board
func (h *HCE) getPhase(b *board.Board) (int, int) {
	phase := -1

	whiteScore, blackScore := 0, 0

	for pce := board.WN; pce <= board.WQ; pce++ {
		whiteScore += b.PieceBitBoards[pce].CountBits() * MaterialScore[Opening][pce]
	}

	for pce := board.BN; pce <= board.BQ; pce++ {
		blackScore += b.PieceBitBoards[pce].CountBits() * -MaterialScore[Opening][pce]
	}

	phaseScore := whiteScore + blackScore

	if phaseScore >= PhaseScore[Opening] {
		phase = Opening
	} else if phaseScore <= PhaseScore[EndGame] {
		phase = EndGame
	} else {
		phase = MiddleGame
	}
	return phase, phaseScore
}
