package board

import (
	"fmt"
)

const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPPRNBQKBNR w KQkq - 0 1"

// parses an FEN string and sets board accordingly
func (b *Board) ParseFEN(fen string) error {

	rank := Eight
	file := A

	strIdx := 0
	fenLen := len(fen)
	piece := Empty

	b.Reset()

	// parse the position of pieces on the board
	var parsePositions = func() error {

		for (rank >= One) && (strIdx < fenLen) {
			count := 1

			switch fen[strIdx] {
			case 'p':
				piece = bP
			case 'r':
				piece = bR
			case 'n':
				piece = bN
			case 'b':
				piece = bB
			case 'q':
				piece = bQ
			case 'k':
				piece = bK
			case 'P':
				piece = wP
			case 'R':
				piece = wR
			case 'N':
				piece = wN
			case 'B':
				piece = wB
			case 'Q':
				piece = wQ
			case 'K':
				piece = wK

			case '1', '2', '3', '4', '5', '6', '7', '8':
				piece = Empty
				count = int(fen[strIdx] - byte('0'))
			case '/', ' ':
				rank--
				file = A
				strIdx++
				continue
			default:
				return fmt.Errorf("failed to parse piece positions")
			}

			for range count {
				sq64 := uint8((rank * 8)) + uint8(file)
				sq120 := Sq64to120[sq64]

				if piece != Empty {
					b.Pieces[sq120] = piece
				}
				file++
			}
			strIdx++

		}
		return nil
	}

	// parse the side to move
	var parseSideToMove = func() error {
		if fen[strIdx] == 'w' {
			b.SideToMove = White
		} else if fen[strIdx] == 'b' {
			b.SideToMove = Black
		} else {
			return fmt.Errorf("failed to parse side to move")
		}
		strIdx += 2
		return nil
	}

	// parse castling permissions
	var parseCastlingPermissions = func() error {

		for range 4 {
			if fen[strIdx] == ' ' {
				break
			}
			switch fen[strIdx] {
			case 'K':
				b.CastlingPermission |= WKSide
			case 'Q':
				b.CastlingPermission |= WQSide
			case 'k':
				b.CastlingPermission |= BKSide
			case 'q':
				b.CastlingPermission |= BQSide
			default:
				return fmt.Errorf("failed to parse castling permissions")
			}
			strIdx++
		}
		strIdx++
		return nil
	}

	// parse en passant square
	var parseEnPassant = func() error {
		if fen[strIdx] != '-' {
			file = File(fen[strIdx] - 'a')
			rank = Rank(fen[strIdx+1] - '0')

			if file < A || file > H || rank < One || rank > Eight {
				return fmt.Errorf("failed to parse enpassant square")
			}

			b.EnPass = FRToSq(file, rank)
		}
		return nil
	}

	if err := parsePositions(); err != nil {
		return fmt.Errorf("failed to parse FEN: %v", err)
	}

	if err := parseSideToMove(); err != nil {
		return fmt.Errorf("failed to parse FEN: %v", err)
	}

	if err := parseCastlingPermissions(); err != nil {
		return fmt.Errorf("failed to parse FEN: %v", err)
	}

	if err := parseEnPassant(); err != nil {
		return fmt.Errorf("failed to parse FEN: %v", err)
	}

	// generate unique position hash
	b.PositionKey = b.GenPositionKey()

	return nil

}
