package main

import (
	"fmt"
	"strconv"
	"strings"
)

func (b *Board) validateFEN(fen string) error {
	elements := strings.Split(fen, " ")

	if len(elements) != 6 {
		return fmt.Errorf("invalid FEN: expected 6 fields separated by spaces, got %d", len(elements))
	}

	ranks := strings.Split(elements[0], "/")

	if len(ranks) != 8 {
		return fmt.Errorf("invalid FEN: expected 8 ranks in piece placement, got %d", len(ranks))
	}

	for _, rank := range ranks {

		if len(rank) == 0 {
			return fmt.Errorf("rank cannot be empty")
		}

		count := 0
		for _, c := range rank {
			if c >= '1' && c <= '8' {
				count += int(c - '0')
			} else if strings.ContainsRune(PieceChar, c) {
				count++
			} else {
				return fmt.Errorf("invalid character in piece placement: %c", c)
			}
		}
		if count != 8 {
			return fmt.Errorf("rank must have exactly 8 squares, got %d", count)
		}
	}

	if elements[1] != "w" && elements[1] != "b" {
		return fmt.Errorf("side to move must be 'w' or 'b'")
	}

	castling := elements[2]
	if castling != "-" {
		for _, c := range castling {
			if !strings.ContainsRune("KQkq", c) {
				return fmt.Errorf("invalid castling character: %c", c)
			}
		}
	}

	ep := elements[3]
	if ep != "-" {
		if len(ep) != 2 || ep[0] < 'a' || ep[0] > 'h' || ep[1] < '1' || ep[1] > '8' {
			return fmt.Errorf("invalid en passant square: %s", ep)
		}
	}

	if _, err := strconv.Atoi(elements[4]); err != nil || elements[4][0] == '-' {
		return fmt.Errorf("invalid halfmove clock: %s", elements[4])
	}

	fullmove, err := strconv.Atoi(elements[5])
	if err != nil || fullmove < 1 {
		return fmt.Errorf("invalid fullmove number: %s", elements[5])
	}

	return nil
}

func (b *Board) ParseFEN(fen string) error {

	if err := b.validateFEN(fen); err != nil {
		return fmt.Errorf("failed to parse fen string: %v", err)
	}

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

	return nil

}
