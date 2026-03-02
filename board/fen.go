package board

import "fmt"

const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPPRNBQKBNR w KQkq - 0 1"

// parses an FEN string and sets board accordingly
func (b *Board) ParseFEN(fen string) error {

	rank := Eight
	file := A

	strIdx := 0
	fenLen := len(fen)
	piece := Empty

	b.Reset()

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
			return fmt.Errorf("failed to parse FEN string")
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
