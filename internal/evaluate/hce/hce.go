package hce

import "ranga/internal/board"

type HCE struct{}

func (h HCE) Evaluate(b *board.Board) int {
	// determines current game phase based on remaining non-pawn material
	phase, phaseScore := h.getPhase(b)

	score := 0
	scoreOpen, scoreEnd := 0, 0

	for pce := board.WP; pce <= board.BK; pce++ {
		bb := b.PieceBitBoards[pce]

		for bb != 0 {
			sq := board.Square(bb.GetLSB())

			// increment static material values
			scoreOpen += MaterialScore[Opening][pce]
			scoreEnd += MaterialScore[EndGame][pce]

			switch pce {
			case board.WP:
				po, pe := h.evalPawns(b, sq, board.White)
				scoreOpen += po
				scoreEnd += pe

			case board.WN:
				no, ne := h.evalKnights(sq, board.White)
				scoreOpen += no
				scoreEnd += ne

			case board.WB:
				bo, be := h.evalBishops(b, sq, board.White)
				scoreOpen += bo
				scoreEnd += be

			case board.WR:
				ro, re := h.evalRooks(b, sq, board.White)
				scoreOpen += ro
				scoreEnd += re

			case board.WQ:
				qo, qe := h.evalQueens(b, sq, board.White)
				scoreOpen += qo
				scoreEnd += qe

			case board.WK:
				ko, ke := h.evalKings(b, sq, board.White)
				scoreOpen += ko
				scoreEnd += ke

			case board.BP:
				po, pe := h.evalPawns(b, sq, board.Black)
				scoreOpen -= po
				scoreEnd -= pe

			case board.BN:
				no, ne := h.evalKnights(sq, board.Black)
				scoreOpen -= no
				scoreEnd -= ne

			case board.BB:
				bo, be := h.evalBishops(b, sq, board.Black)
				scoreOpen -= bo
				scoreEnd -= be

			case board.BR:
				ro, re := h.evalRooks(b, sq, board.Black)
				scoreOpen -= ro
				scoreEnd -= re

			case board.BQ:
				qo, qe := h.evalQueens(b, sq, board.Black)
				scoreOpen -= qo
				scoreEnd -= qe

			case board.BK:
				ko, ke := h.evalKings(b, sq, board.Black)
				scoreOpen -= ko
				scoreEnd -= ke
			}

			bb.PopBit(sq)

		}
	}

	// taper score between opening and endgame phases
	score = h.taper(scoreOpen, scoreEnd, phase, phaseScore)

	if b.Side == board.Black {
		return -score
	}

	return score

}
