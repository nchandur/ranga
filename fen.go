package main

func (b *Board) ParseFEN(fen string) {

	b.Reset()

	rank := Rank8
	file := FileA
	piece := Empty
	sq120 := Square(0)

	fenCount := 0
	n := len(fen)

	for rank >= Rank1 && fenCount < n {
		count := 1

		switch fen[fenCount] {
		case 'p':
			piece = bP
		case 'n':
			piece = bN
		case 'b':
			piece = bB
		case 'r':
			piece = bR
		case 'q':
			piece = bQ
		case 'k':
			piece = bK

		case 'P':
			piece = wP
		case 'N':
			piece = wN
		case 'B':
			piece = wB
		case 'R':
			piece = wR
		case 'Q':
			piece = wQ
		case 'K':
			piece = wK

		case '1', '2', '3', '4', '5', '6', '7', '8':
			piece = Empty
			count = int(fen[fenCount] - '0')

		case '/', ' ':
			rank--
			file = FileA
			fenCount++
			continue
		}

		for range count {
			sq120 = FR2Square(file, rank)
			b.Pieces[sq120] = piece
			file++
		}

		fenCount++

	}

	switch fen[fenCount] {
	case 'w':
		b.SideToMove = White
	case 'b':
		b.SideToMove = Black
	default:
		b.SideToMove = Both
	}

	fenCount += 2

	for range 4 {
		if fen[fenCount] == ' ' {
			break
		}

		switch fen[fenCount] {
		case 'K':
			b.Castling |= WKSide
		case 'Q':
			b.Castling |= WQSide
		case 'k':
			b.Castling |= BKSide
		case 'q':
			b.Castling |= BQSide
		}

		fenCount++

	}

	fenCount++

	if fen[fenCount] == '-' {
		b.EnPassant = NoSquare
	} else {
		file := File(fen[fenCount] - 'a')
		rank := Rank(fen[fenCount+1] - '1')
		b.EnPassant = FR2Square(file, rank)
	}

	b.Hash = b.GenHash()
	b.UpdatePieceList()

}
