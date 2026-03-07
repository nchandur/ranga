package board

import (
	"fmt"
	"strings"
)

const maxMoves = 256

type MoveList struct {
	moves []Move
	count int
}

func NewMoveList() MoveList {
	ml := MoveList{}
	ml.moves = make([]Move, maxMoves)
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

func (m *MoveList) addEnPassantMove(move int) {
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

func (m *MoveList) Generate(board *Board) {

	if !board.Check() {
		return
	}

	m.pawnMoveGeneration(board)

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
