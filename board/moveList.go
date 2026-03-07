package board

import (
	"fmt"
	"strings"
)

var pieceDirection = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
}

// number of times to loop through in that direction
var numDirection = []int{
	0, 0, 8, 4, 4, 8, 8, 0, 8, 4, 4, 8, 8,
}

type MoveList struct {
	moves []Move
	count int
}

func NewMoveList() MoveList {
	ml := MoveList{}
	ml.moves = make([]Move, 256)
	return ml
}

func (m *MoveList) addQuietMove(move int) {
	m.moves[m.count].move = move
	m.moves[m.count].score = 0
	m.count++
}

func (m *MoveList) addCaptureMove(move int) {
	m.moves[m.count].move = move
	m.moves[m.count].score = 0
	m.count++
}

func (m *MoveList) addWhitePawnCaptureMove(from, to Square, capture Piece) {
	if Fr120ToRank(int(from)) == Seven {
		m.addCaptureMove(NewMove(from, to, capture, wQ, false, false).move)
		m.addCaptureMove(NewMove(from, to, capture, wR, false, false).move)
		m.addCaptureMove(NewMove(from, to, capture, wB, false, false).move)
		m.addCaptureMove(NewMove(from, to, capture, wK, false, false).move)
	} else {
		m.addCaptureMove(NewMove(from, to, capture, Empty, false, false).move)
	}
}

func (m *MoveList) addWhitePawnMove(from, to Square) {
	if Fr120ToRank(int(from)) == Seven {
		m.addQuietMove(NewMove(from, to, Empty, wQ, false, false).move)
		m.addQuietMove(NewMove(from, to, Empty, wR, false, false).move)
		m.addQuietMove(NewMove(from, to, Empty, wB, false, false).move)
		m.addQuietMove(NewMove(from, to, Empty, wK, false, false).move)
	} else {
		m.addQuietMove(NewMove(from, to, Empty, Empty, false, false).move)
	}
}

func (m *MoveList) addBlackPawnCaptureMove(from, to Square, capture Piece) {
	if Fr120ToRank(int(from)) == Two {
		m.addCaptureMove(NewMove(from, to, capture, bQ, false, false).move)
		m.addCaptureMove(NewMove(from, to, capture, bR, false, false).move)
		m.addCaptureMove(NewMove(from, to, capture, bB, false, false).move)
		m.addCaptureMove(NewMove(from, to, capture, bK, false, false).move)
	} else {
		m.addCaptureMove(NewMove(from, to, capture, Empty, false, false).move)
	}
}

func (m *MoveList) addBlackPawnMove(from, to Square) {
	if Fr120ToRank(int(from)) == Two {
		m.addQuietMove(NewMove(from, to, Empty, bQ, false, false).move)
		m.addQuietMove(NewMove(from, to, Empty, bR, false, false).move)
		m.addQuietMove(NewMove(from, to, Empty, bB, false, false).move)
		m.addQuietMove(NewMove(from, to, Empty, bK, false, false).move)
	} else {
		m.addQuietMove(NewMove(from, to, Empty, Empty, false, false).move)
	}
}

func (m *MoveList) pawnMoveGeneration(board *Board) {
	m.count = 0
	side := board.SideToMove

	switch side {
	case White:
		for pieceNum := range board.PieceNumber[wP] {

			sq := Square(board.PieceList[wP][pieceNum])

			if board.Pieces[sq+10] == Empty {
				m.addWhitePawnMove(sq, sq+10)
				if Fr120ToRank(int(sq)) == Two && board.Pieces[sq+20] == Empty {
					m.addQuietMove(NewMove(sq, sq+20, Empty, Empty, false, true).move)
				}

			}

			tempSq := board.Pieces[sq+9]
			if tempSq != 120 && pieceColor[tempSq] == Black {
				m.addWhitePawnCaptureMove(sq, sq+9, board.Pieces[sq+9])
			}

			tempSq = board.Pieces[sq+11]
			if tempSq != 120 && pieceColor[tempSq] == Black {
				m.addWhitePawnCaptureMove(sq, sq+11, board.Pieces[sq+11])
			}

			if (sq + 9) == board.EnPass {
				m.addCaptureMove(NewMove(sq, sq+9, Empty, Empty, true, false).move)
			}

			if (sq + 11) == board.EnPass {
				m.addCaptureMove(NewMove(sq, sq+11, Empty, Empty, true, false).move)
			}

		}
	case Black:
		for pieceNum := range board.PieceNumber[bP] {

			sq := Square(board.PieceList[bP][pieceNum])

			if board.Pieces[sq-10] == Empty {
				m.addBlackPawnMove(sq, sq-10)
				if Fr120ToRank(int(sq)) == Seven && board.Pieces[sq-20] == Empty {
					m.addQuietMove(NewMove(sq, sq-20, Empty, Empty, false, true).move)
				}

			}

			tempSq := board.Pieces[sq-9]
			if tempSq != 120 && pieceColor[tempSq] == White {
				m.addBlackPawnCaptureMove(sq, sq-9, board.Pieces[sq-9])
			}

			tempSq = board.Pieces[sq-11]
			if tempSq != 120 && pieceColor[tempSq] == White {
				m.addBlackPawnCaptureMove(sq, sq-11, board.Pieces[sq-11])
			}

			if (sq - 9) == board.EnPass {
				m.addCaptureMove(NewMove(sq, sq-9, Empty, Empty, true, false).move)
			}

			if (sq - 11) == board.EnPass {
				m.addCaptureMove(NewMove(sq, sq-11, Empty, Empty, true, false).move)
			}

		}

	}

}

func (m *MoveList) slidingPieceMoveGeneration(board *Board) {
	slidePiece := []Piece{wB, wR, wQ, Empty, bB, bR, bQ, Empty}
	slideIdx := []int{0, 4}

	pieceIdx := slideIdx[board.SideToMove]
	piece := slidePiece[pieceIdx]
	pieceIdx++

	for piece != Empty {

		for pieceNum := range board.PieceNumber[piece] {
			sq := Square(board.PieceList[piece][pieceNum])

			for idx := range numDirection[piece] {
				dir := pieceDirection[piece][idx]
				tempSq := Square(int(sq) + dir)

				// if square is offboard
				for board.Pieces[tempSq] != 120 {
					// if square is a capture
					if board.Pieces[tempSq] != Empty {
						if pieceColor[board.Pieces[tempSq]] == board.SideToMove^1 {
							m.addCaptureMove(NewMove(sq, tempSq, board.Pieces[tempSq], Empty, false, false).move)
						}
						break
					}

					m.addQuietMove(NewMove(sq, tempSq, Empty, Empty, false, false).move)
					tempSq = Square(int(tempSq) + dir)
				}

			}

		}

		piece = slidePiece[pieceIdx]
		pieceIdx++
	}

}

func (m *MoveList) nonSlidingPieceMoveGeneration(board *Board) {
	nonSlidePiece := []Piece{wN, wK, Empty, bN, bK, Empty}
	nonSlideIdx := []int{0, 3}

	pieceIdx := nonSlideIdx[board.SideToMove]
	piece := nonSlidePiece[pieceIdx]
	pieceIdx++

	for piece != Empty {

		for pieceNum := range board.PieceNumber[piece] {
			sq := Square(board.PieceList[piece][pieceNum])
			for idx := range numDirection[piece] {
				dir := pieceDirection[piece][idx]
				tempSq := Square(int(sq) + dir)

				// if square is offboard
				if board.Pieces[tempSq] == 120 {
					continue
				}

				// if square is a capture
				if board.Pieces[tempSq] != Empty {
					if pieceColor[board.Pieces[tempSq]] == board.SideToMove^1 {
						m.addCaptureMove(NewMove(sq, tempSq, board.Pieces[tempSq], Empty, false, false).move)
					}
					continue
				}

				m.addQuietMove(NewMove(sq, tempSq, Empty, Empty, false, false).move)

			}

		}

		piece = nonSlidePiece[pieceIdx]
		pieceIdx++

	}

}

func (m *MoveList) castlingMoveGeneration(board *Board) {

	switch board.SideToMove {
	case White:
		if board.CastlingPermission&WKSide != 0 {
			if board.Pieces[F1] == Empty && board.Pieces[G1] == Empty {
				if !board.IsAttacked(E1, Black) && !board.IsAttacked(F1, Black) {
					m.addQuietMove(NewMove(E1, G1, Empty, Empty, false, false).move)
				}
			}
		}

		if board.CastlingPermission&WQSide != 0 {
			if board.Pieces[D1] == Empty && board.Pieces[C1] == Empty && board.Pieces[B1] == Empty {
				if !board.IsAttacked(E1, Black) && !board.IsAttacked(D1, Black) {
					m.addQuietMove(NewMove(E1, C1, Empty, Empty, false, false).move)
				}
			}
		}

	case Black:

		if board.CastlingPermission&BKSide != 0 {
			if board.Pieces[F8] == Empty && board.Pieces[G8] == Empty {
				if !board.IsAttacked(E8, White) && !board.IsAttacked(F8, White) {
					m.addQuietMove(NewMove(E8, G8, Empty, Empty, false, false).move)
				}
			}
		}

		if board.CastlingPermission&BQSide != 0 {
			if board.Pieces[D8] == Empty && board.Pieces[C8] == Empty && board.Pieces[B8] == Empty {
				if !board.IsAttacked(E8, White) && !board.IsAttacked(D8, White) {
					m.addQuietMove(NewMove(E8, C8, Empty, Empty, false, false).move)
				}
			}
		}

	}

}

func (m *MoveList) Generate(board *Board) {

	if !board.Check() {
		return
	}

	m.pawnMoveGeneration(board)
	m.slidingPieceMoveGeneration(board)
	m.nonSlidingPieceMoveGeneration(board)
	m.castlingMoveGeneration(board)

}

func (m *MoveList) String() string {
	var builder strings.Builder

	for idx := range m.count {
		score := m.moves[idx].score
		fmt.Fprintf(&builder, "Move #%d: %s (score: %d)\n", idx+1, m.moves[idx].String(), score)
	}

	fmt.Fprintf(&builder, "Move List: %d\n", m.count)
	return builder.String()
}
