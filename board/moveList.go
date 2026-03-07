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

func (m *MoveList) AddQuietMove(board *Board, move int) {
	m.moves[m.count].move = move
	m.moves[m.count].score = 0
	m.count++
}

func (m *MoveList) AddCaptureMove(board *Board, move int) {
	m.moves[m.count].move = move
	m.moves[m.count].score = 0
	m.count++
}

func (m *MoveList) AddEnPassantMove(board *Board, move int) {
	m.moves[m.count].move = move
	m.moves[m.count].score = 0
	m.count++
}

func (m *MoveList) AddWhitePawnCaptureMove(board *Board, from, to Square, capture Piece) {
	if Fr120ToRank(int(from)) == Seven {
		m.AddCaptureMove(board, NewMove(from, to, capture, wQ, false, false).move)
		m.AddCaptureMove(board, NewMove(from, to, capture, wR, false, false).move)
		m.AddCaptureMove(board, NewMove(from, to, capture, wB, false, false).move)
		m.AddCaptureMove(board, NewMove(from, to, capture, wK, false, false).move)
	} else {
		m.AddCaptureMove(board, NewMove(from, to, capture, Empty, false, false).move)
	}
}

func (m *MoveList) AddWhitePawnMove(board *Board, from, to Square) {
	if Fr120ToRank(int(from)) == Seven {
		m.AddQuietMove(board, NewMove(from, to, Empty, wQ, false, false).move)
		m.AddQuietMove(board, NewMove(from, to, Empty, wR, false, false).move)
		m.AddQuietMove(board, NewMove(from, to, Empty, wB, false, false).move)
		m.AddQuietMove(board, NewMove(from, to, Empty, wK, false, false).move)
	} else {
		m.AddQuietMove(board, NewMove(from, to, Empty, Empty, false, false).move)
	}
}

func (m *MoveList) Generate(board *Board) {

	if !board.Check() {
		return
	}

	m.count = 0
	side := board.SideToMove

	if side == White {
		for pieceNum := range board.PieceNumber[wP] {

			sq := Square(board.PieceList[wP][pieceNum])

			if board.Pieces[sq+10] == Empty {
				m.AddWhitePawnMove(board, sq, sq+10)
				if Fr120ToRank(int(sq)) == Two && board.Pieces[sq+20] == Empty {
					m.AddQuietMove(board, NewMove(sq, sq+20, Empty, Empty, false, true).move)
				}

			}

			tempSq := board.Pieces[sq+9]
			if tempSq != 120 && pieceColor[tempSq] == Black {
				m.AddWhitePawnCaptureMove(board, sq, sq+9, board.Pieces[sq+9])
			}

			tempSq = board.Pieces[sq+11]
			if tempSq != 120 && pieceColor[tempSq] == Black {
				m.AddWhitePawnCaptureMove(board, sq, sq+11, board.Pieces[sq+11])
			}

			if (sq + 9) == board.EnPass {
				m.AddCaptureMove(board, NewMove(sq, sq+9, Empty, Empty, true, false).move)
			}

			if (sq + 11) == board.EnPass {
				m.AddCaptureMove(board, NewMove(sq, sq+11, Empty, Empty, true, false).move)
			}

		}
	}

	if side == Black {

	}

}

func (m *MoveList) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Move List: %d\n", m.count)

	for idx := range m.count {
		score := m.moves[idx].score
		fmt.Fprintf(&builder, "Move #%d: %s Promoted: %d\t(score: %d)\n", idx+1, m.moves[idx].String(), m.moves[idx].PromotedPiece(), score)
	}

	return builder.String()
}
