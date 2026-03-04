package board

import (
	"fmt"
	"strings"
)

// parses an FEN string and sets board accordingly
func (b *Board) ParseFEN(fen string) error {

	fenSlice := strings.Split(fen, " ")

	rank := Eight
	file := A

	b.Reset()

	var parsePositions = func(positions string) error {

		rank := Eight
		file := One

		for i := range positions {

			count := 1
			piece := Empty

			switch positions[i] {

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
				count = int(positions[i] - '0')

			case '/':
				rank--
				file = 0
				continue

			default:
				return fmt.Errorf("failed to parse piece positions: %c", positions[i])
			}

			for j := 0; j < count; j++ {

				if rank < 0 || rank > 7 || file < 0 || file > 7 {
					return fmt.Errorf("rank/file out of bounds r=%d f=%d", rank, file)
				}

				sq120 := FRToSq(File(file), Rank(rank))

				if piece != Empty {
					b.Pieces[sq120] = piece
				}

				file++
			}
		}

		return nil
	}
	var parseSideToMove = func(side string) error {

		switch side {
		case "w":
			b.SideToMove = White
		case "b":
			b.SideToMove = Black
		default:
			return fmt.Errorf("failed to parse side to move")
		}

		return nil
	}

	var parseCastlingPermissions = func(castle string) error {

		// no castling rights
		if castle == "-" {
			return nil
		}

		// check for castling rights on K and Q sides
		for _, c := range castle {
			switch c {
			case 'K':
				b.CastlingPermission |= WKSide
			case 'Q':
				b.CastlingPermission |= WQSide
			case 'k':
				b.CastlingPermission |= BKSide
			case 'q':
				b.CastlingPermission |= BQSide
			default:
				return fmt.Errorf("failed to parse castling permissions: %c", c)
			}
		}

		return nil
	}

	// parse en passant square
	var parseEnPassant = func(enpassant string) error {

		// no en passant square
		if enpassant == "-" {
			return nil
		}

		if len(enpassant) > 2 {
			return fmt.Errorf("failed to parse enpassant square")
		}

		file = File(enpassant[0] - 'a')
		rank = Rank(enpassant[1] - '1')

		if file < A || file > H || rank < One || rank > Eight {
			return fmt.Errorf("failed to parse enpassant square")
		}

		b.EnPass = FRToSq(file, rank)
		return nil
	}

	if err := parsePositions(fenSlice[0]); err != nil {
		return err
	}

	if err := parseSideToMove(fenSlice[1]); err != nil {
		return err
	}

	if err := parseCastlingPermissions(fenSlice[2]); err != nil {
		return err
	}

	if err := parseEnPassant(fenSlice[3]); err != nil {
		return err
	}

	b.PositionKey = b.GenPositionKey()

	return nil
}
