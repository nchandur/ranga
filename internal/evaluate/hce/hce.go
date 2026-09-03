package hce

import "ranga/internal/board"

type HCE struct{}

func (h HCE) Evaluate(b *board.Board) int {
	score := 0

	for pce := board.WP; pce <= board.BK; pce++ {
		bb := b.PieceBitBoards[pce]

		for bb != 0 {
			sq := board.Square(bb.GetLSB())

			switch pce {
			case board.WP:
				score += PawnTable[sq]

				doublePawn := (b.PieceBitBoards[pce] & FileMasks[sq]).CountBits()

				// double pawn
				if doublePawn > 1 {
					score += doublePawn * doublePawnPenalty
				}

				// isolated pawn
				if (b.PieceBitBoards[pce] & IsolatedMasks[sq]) == 0 {
					score += isolatedPawnPenalty
				}

				// passed pawn
				if (WhitePassedPawnMasks[sq] & b.PieceBitBoards[board.BP]) == 0 {
					score += PassedPawnBonus[Ranks[sq]]
				}
			case board.WN:
				score += KnightTable[sq]
			case board.WB:
				score += BishopTable[sq]
				score += board.GetBishopAttacks(sq, b.Occupancies[board.Both]).CountBits()
			case board.WR:
				score += RookTable[sq]
				score += board.GetRookAttacks(sq, b.Occupancies[board.Both]).CountBits()

				// semi open file
				if (b.PieceBitBoards[board.WP] & FileMasks[sq]) == 0 {
					score += semiOpenFile
				}
				// open file
				if ((b.PieceBitBoards[board.WP] | b.PieceBitBoards[board.BP]) & FileMasks[sq]) == 0 {
					score += openFile
				}
			case board.WQ:
				score += board.GetQueenAttacks(sq, b.Occupancies[board.Both]).CountBits()
			case board.WK:
				score += KingTable[sq]

				// semi open file
				if (b.PieceBitBoards[board.WP] & FileMasks[sq]) == 0 {
					score -= semiOpenFile
				}
				// open file
				if ((b.PieceBitBoards[board.WP] | b.PieceBitBoards[board.BP]) & FileMasks[sq]) == 0 {
					score -= openFile
				}

				// king safety
				score += ((board.KingAttacks[sq] & b.Occupancies[board.White]).CountBits() * kingSafetyBonus)

			case board.BP:
				score -= PawnTable[sq^56]

				doublePawns := (b.PieceBitBoards[pce] & FileMasks[sq]).CountBits()

				// double pawns
				if doublePawns > 1 {
					score -= doublePawns * doublePawnPenalty
				}

				// isolated pawns
				if (b.PieceBitBoards[pce] & IsolatedMasks[sq]) == 0 {
					score -= isolatedPawnPenalty
				}

				// passed pawns
				if (BlackPassedPawnMasks[sq] & b.PieceBitBoards[board.WP]) == 0 {
					score -= PassedPawnBonus[Ranks[sq^56]]
				}

			case board.BN:
				score -= KnightTable[sq^56]
			case board.BB:
				score -= BishopTable[sq^56]
				score -= board.GetBishopAttacks(sq, b.Occupancies[board.Both]).CountBits()

			case board.BR:
				score -= RookTable[sq^56]

				// semi open file
				if (b.PieceBitBoards[board.BP] & FileMasks[sq]) == 0 {
					score -= semiOpenFile
				}
				// open file
				if ((b.PieceBitBoards[board.WP] | b.PieceBitBoards[board.BP]) & FileMasks[sq]) == 0 {
					score -= openFile
				}

			case board.BQ:
				score -= board.GetQueenAttacks(sq, b.Occupancies[board.Both]).CountBits()
			case board.BK:
				score -= KingTable[sq^56]

				// semi open file
				if (b.PieceBitBoards[board.BP] & FileMasks[sq]) == 0 {
					score += semiOpenFile
				}
				// open file
				if ((b.PieceBitBoards[board.WP] | b.PieceBitBoards[board.BP]) & FileMasks[sq]) == 0 {
					score += openFile
				}

				// king safety
				score -= ((board.KingAttacks[sq] & b.Occupancies[board.Black]).CountBits() * kingSafetyBonus)

			}

			score += board.PieceValue[pce]
			bb.PopBit(sq)

		}
	}

	if b.Side == board.Black {
		return -score
	}

	return score

}
