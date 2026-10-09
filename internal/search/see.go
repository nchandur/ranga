package search

import "ranga/internal/board"

// SeeCapture calculates the exact material exchange of a specific capture move.
func (s *Searcher) SeeCapture(b *board.Board, move board.Move) int {
	from := move.Source()
	to := move.Target()
	piece := move.Piece()

	targetPiece := b.Mailbox[to]
	capturedVal := 0
	if targetPiece != board.Empty {
		capturedVal = abs(board.PieceValue[targetPiece])
	}

	occupied := b.Occupancies[board.Both]

	// Handle En Passant captures by removing the invisible target pawn
	if move.IsEnpass() {
		capturedVal = abs(board.PieceValue[board.WP])
		if b.Side == board.White {
			occupied.PopBit(to + 8)
		} else {
			occupied.PopBit(to - 8)
		}
	}

	// Simulated 'make_capture' against the starting piece
	occupied.PopBit(from)

	return capturedVal - s.see(b, to, piece, b.Side^1, occupied)
}

// see executes the recursive evaluation of captures on a single square.
func (s *Searcher) see(b *board.Board, sq board.Square, pieceOnSquare board.Piece, sideToMove board.Color, occupied board.BitBoard) int {
	value := 0

	attacker, fromSq := s.getSmallestAttacker(b, sq, sideToMove, occupied)

	// skip if the square isn't attacked anymore by this side
	if attacker != board.Empty {
		// Simulated 'make_capture' to expose sliding pieces (x-rays) behind the attacker
		occupied.PopBit(fromSq)

		// The piece just captured is the unit sitting on the square during this iteration
		capturedVal := abs(board.PieceValue[pieceOnSquare])

		// max(0, piece_just_captured - see(square, other_side))
		score := capturedVal - s.see(b, sq, attacker, sideToMove^1, occupied)
		if score > 0 {
			value = score
		}
	}

	return value
}

// getSmallestAttacker sequentially searches for the least valuable attacker of a square.
func (s *Searcher) getSmallestAttacker(b *board.Board, sq board.Square, side board.Color, occupied board.BitBoard) (board.Piece, board.Square) {
	if side == board.White {
		// Evaluate reverse attacks for pawns to detect valid attackers
		if pawns := board.MaskPawnAttacks(sq, board.Black) & b.PieceBitBoards[board.WP] & occupied; pawns != 0 {
			return board.WP, board.Square(pawns.GetLSB())
		}
		if knights := board.MaskKnightAttacks(sq) & b.PieceBitBoards[board.WN] & occupied; knights != 0 {
			return board.WN, board.Square(knights.GetLSB())
		}
		if bishops := board.GetBishopAttacks(sq, occupied) & b.PieceBitBoards[board.WB] & occupied; bishops != 0 {
			return board.WB, board.Square(bishops.GetLSB())
		}
		if rooks := board.GetRookAttacks(sq, occupied) & b.PieceBitBoards[board.WR] & occupied; rooks != 0 {
			return board.WR, board.Square(rooks.GetLSB())
		}
		if queens := board.GetQueenAttacks(sq, occupied) & b.PieceBitBoards[board.WQ] & occupied; queens != 0 {
			return board.WQ, board.Square(queens.GetLSB())
		}
		if kings := board.MaskKingAttacks(sq) & b.PieceBitBoards[board.WK] & occupied; kings != 0 {
			return board.WK, board.Square(kings.GetLSB())
		}
	} else {
		if pawns := board.MaskPawnAttacks(sq, board.White) & b.PieceBitBoards[board.BP] & occupied; pawns != 0 {
			return board.BP, board.Square(pawns.GetLSB())
		}
		if knights := board.MaskKnightAttacks(sq) & b.PieceBitBoards[board.BN] & occupied; knights != 0 {
			return board.BN, board.Square(knights.GetLSB())
		}
		if bishops := board.GetBishopAttacks(sq, occupied) & b.PieceBitBoards[board.BB] & occupied; bishops != 0 {
			return board.BB, board.Square(bishops.GetLSB())
		}
		if rooks := board.GetRookAttacks(sq, occupied) & b.PieceBitBoards[board.BR] & occupied; rooks != 0 {
			return board.BR, board.Square(rooks.GetLSB())
		}
		if queens := board.GetQueenAttacks(sq, occupied) & b.PieceBitBoards[board.BQ] & occupied; queens != 0 {
			return board.BQ, board.Square(queens.GetLSB())
		}
		if kings := board.MaskKingAttacks(sq) & b.PieceBitBoards[board.BK] & occupied; kings != 0 {
			return board.BK, board.Square(kings.GetLSB())
		}
	}

	return board.Empty, board.NoSquare
}
