package main

func (b *Board) addCaptureMove(move Move) {
	b.MoveList[b.MoveListStart[b.Ply+1]] = int(move)
	b.MoveScores[b.MoveListStart[b.Ply+1]] = 0
	b.MoveListStart[b.Ply+1]++
}

func (b *Board) addQuietMove(move Move) {
	b.MoveList[b.MoveListStart[b.Ply+1]] = int(move)
	b.MoveScores[b.MoveListStart[b.Ply+1]] = 0
	b.MoveListStart[b.Ply+1]++
}

func (b *Board) addEnpassantMove(move Move) {
	b.MoveList[b.MoveListStart[b.Ply+1]] = int(move)
	b.MoveScores[b.MoveListStart[b.Ply+1]] = 0
	b.MoveListStart[b.Ply+1]++
}

func (b *Board) genNonSlideMoves() {
	pieceIdx := nonSlidePieceIdx[b.SideToMove]
	piece := nonSlidePiece[pieceIdx]
	pieceIdx++

	for piece != 0 {

		for i := 0; i < b.PieceNumber[piece]; i++ {

			sq := b.PieceList[(int(piece)*10)+i]

			for idx := 0; idx < dirNum[piece]; idx++ {

				dir := pieceDir[piece][idx]
				tSq := sq + dir

				if FilesBoard[tSq] == int(Offboard) {
					continue
				}

				if b.Pieces[tSq] != Empty {

					if pieceColor[b.Pieces[tSq]] != b.SideToMove {
						b.addCaptureMove(NewMove(sq, tSq, b.Pieces[tSq], Empty, MFLAGCAP))
					}

				} else {
					b.addQuietMove(NewMove(sq, tSq, Empty, Empty, MNONE))
				}
			}
		}

		piece = nonSlidePiece[pieceIdx]
		pieceIdx++
	}
}

func (b *Board) genSlideMoves() {
	pieceIdx := slidePieceIdx[b.SideToMove]
	piece := slidePiece[pieceIdx]
	pieceIdx++

	for piece != 0 {

		for i := 0; i < b.PieceNumber[piece]; i++ {

			sq := b.PieceList[(int(piece)*10)+i]

			for idx := 0; idx < dirNum[piece]; idx++ {

				dir := pieceDir[piece][idx]
				tSq := sq + dir

				for FilesBoard[tSq] != int(Offboard) {

					if b.Pieces[tSq] != Empty {

						if pieceColor[b.Pieces[tSq]] != b.SideToMove {
							b.addCaptureMove(NewMove(sq, tSq, b.Pieces[tSq], Empty, MFLAGCAP))
						}

						break
					}

					b.addQuietMove(NewMove(sq, tSq, Empty, Empty, MNONE))

					tSq += dir
				}
			}
		}

		piece = slidePiece[pieceIdx]
		pieceIdx++
	}
}

func (b *Board) genPawnCaptureMove(from, to Square, captured Piece, side Color) {
	switch side {
	case White:
		if RanksBoard[from] == int(Rank7) {
			b.addCaptureMove(NewMove(from, to, captured, wN, MNONE))
			b.addCaptureMove(NewMove(from, to, captured, wB, MNONE))
			b.addCaptureMove(NewMove(from, to, captured, wR, MNONE))
			b.addCaptureMove(NewMove(from, to, captured, wQ, MNONE))
		} else {
			b.addCaptureMove(NewMove(from, to, captured, Empty, MFLAGCAP))
		}

	case Black:
		if RanksBoard[from] == int(Rank2) {
			b.addCaptureMove(NewMove(from, to, captured, bN, MNONE))
			b.addCaptureMove(NewMove(from, to, captured, bB, MNONE))
			b.addCaptureMove(NewMove(from, to, captured, bR, MNONE))
			b.addCaptureMove(NewMove(from, to, captured, bQ, MNONE))
		} else {
			b.addCaptureMove(NewMove(from, to, captured, Empty, MFLAGCAP))
		}
	}

}

func (b *Board) genPawnQuietMove(from, to Square, side Color) {
	switch side {
	case White:
		if RanksBoard[from] == int(Rank7) {
			b.addQuietMove(NewMove(from, to, Empty, wN, MNONE))
			b.addQuietMove(NewMove(from, to, Empty, wB, MNONE))
			b.addQuietMove(NewMove(from, to, Empty, wR, MNONE))
			b.addQuietMove(NewMove(from, to, Empty, wQ, MNONE))
		} else {
			b.addQuietMove(NewMove(from, to, Empty, Empty, MNONE))
		}

	case Black:
		if RanksBoard[from] == int(Rank2) {
			b.addQuietMove(NewMove(from, to, Empty, bN, MNONE))
			b.addQuietMove(NewMove(from, to, Empty, bB, MNONE))
			b.addQuietMove(NewMove(from, to, Empty, bR, MNONE))
			b.addQuietMove(NewMove(from, to, Empty, bQ, MNONE))
		} else {
			b.addQuietMove(NewMove(from, to, Empty, Empty, MNONE))
		}
	}

}

func (b *Board) GenMoves() {

	b.MoveListStart[b.Ply+1] = b.MoveListStart[b.Ply]

	switch b.SideToMove {
	case White:

		for idx := range b.PieceNumber[wP] {
			sq := b.PieceList[(int(wP)*10)+idx]

			if b.Pieces[sq+10] == Empty {
				b.genPawnQuietMove(sq, sq+10, White)
				if RanksBoard[sq] == int(Rank2) && b.Pieces[sq+20] == Empty {
					b.addQuietMove(NewMove(sq, sq+20, Empty, Empty, MFLAGPS))
				}
			}

			if FilesBoard[sq+9] != int(Offboard) && pieceColor[b.Pieces[sq+9]] == Black {
				b.genPawnCaptureMove(sq, sq+9, b.Pieces[sq+9], White)
			}

			if FilesBoard[sq+11] != int(Offboard) && pieceColor[b.Pieces[sq+11]] == Black {
				b.genPawnCaptureMove(sq, sq+11, b.Pieces[sq+11], White)
			}

			if b.EnPassant != NoSquare {
				if sq+9 == b.EnPassant {
					b.addEnpassantMove(NewMove(sq, sq+9, Empty, Empty, MFLAGEP))
				}

				if sq+11 == b.EnPassant {
					b.addEnpassantMove(NewMove(sq, sq+11, Empty, Empty, MFLAGEP))
				}
			}

		}

		if b.Castling&WKSide != 0 {
			if (b.Pieces[F1] == Empty) && b.Pieces[G1] == Empty {
				if !b.IsAttacked(F1, Black) && !b.IsAttacked(E1, Black) {
					b.addQuietMove(NewMove(E1, G1, Empty, Empty, MFLAGCA))
				}
			}
		}

		if b.Castling&WQSide != 0 {
			if b.Pieces[D1] == Empty && b.Pieces[C1] == Empty && b.Pieces[B1] == Empty {
				if !b.IsAttacked(D1, Black) && !b.IsAttacked(E1, Black) {
					b.addQuietMove(NewMove(E1, C1, Empty, Empty, MFLAGCA))
				}
			}
		}

	case Black:
		for idx := range b.PieceNumber[bP] {
			sq := b.PieceList[(int(bP)*10)+idx]

			if b.Pieces[sq-10] == Empty {
				b.genPawnQuietMove(sq, sq-10, Black)
				if RanksBoard[sq] == int(Rank7) && b.Pieces[sq-20] == Empty {
					b.addQuietMove(NewMove(sq, sq-20, Empty, Empty, MFLAGPS))
				}
			}

			if FilesBoard[sq-9] != int(Offboard) && pieceColor[b.Pieces[sq-9]] == White {
				b.genPawnCaptureMove(sq, sq-9, b.Pieces[sq-9], Black)
			}

			if FilesBoard[sq-11] != int(Offboard) && pieceColor[b.Pieces[sq-11]] == White {
				b.genPawnCaptureMove(sq, sq-11, b.Pieces[sq-11], Black)
			}

			if b.EnPassant != NoSquare {
				if sq-9 == b.EnPassant {
					b.addEnpassantMove(NewMove(sq, sq-9, Empty, Empty, MFLAGEP))
				}

				if sq-11 == b.EnPassant {
					b.addEnpassantMove(NewMove(sq, sq-11, Empty, Empty, MFLAGEP))
				}
			}

		}
		if b.Castling&BKSide != 0 {
			if (b.Pieces[F8] == Empty) && b.Pieces[G8] == Empty {
				if !b.IsAttacked(F8, White) && !b.IsAttacked(E8, White) {
					b.addQuietMove(NewMove(E8, G8, Empty, Empty, MFLAGCA))
				}
			}
		}

		if b.Castling&BQSide != 0 {
			if b.Pieces[D8] == Empty && b.Pieces[C8] == Empty && b.Pieces[B8] == Empty {
				if !b.IsAttacked(D8, White) && !b.IsAttacked(E8, White) {
					b.addQuietMove(NewMove(E8, C8, Empty, Empty, MFLAGCA))
				}
			}
		}

	}

	b.genNonSlideMoves()
	b.genSlideMoves()
}
