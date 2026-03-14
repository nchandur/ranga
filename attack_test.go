package main

import (
	"testing"
)

func TestIsSquareAttackedByPawn(t *testing.T) {
	board := NewBoard()

	board.Pieces[54] = wP
	board.Pieces[65] = bP

	tests := []struct {
		square   Square
		side     Color
		expected bool
	}{
		{square: 65, side: White, expected: true},
		{square: 76, side: White, expected: false},
		{square: 54, side: Black, expected: true},
		{square: 43, side: Black, expected: false},
	}

	for _, test := range tests {
		result := board.isAttackedByPawn(test.square, test.side)
		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}

func TestIsSquareAttackedByKnight(t *testing.T) {
	board := NewBoard()

	board.Pieces[55] = wN
	board.Pieces[92] = bN

	tests := []struct {
		square   Square
		side     Color
		expected bool
	}{
		{square: 34, side: White, expected: true},
		{square: 36, side: White, expected: true},
		{square: 43, side: White, expected: true},
		{square: 47, side: White, expected: true},
		{square: 63, side: White, expected: true},
		{square: 67, side: White, expected: true},
		{square: 74, side: White, expected: true},
		{square: 76, side: White, expected: true},
		{square: 35, side: White, expected: false},
		{square: 111, side: Black, expected: false},
		{square: 113, side: Black, expected: false},
	}

	for _, test := range tests {
		result := board.isAttackedByKnight(test.square, test.side)
		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}

func TestIsSquareAttackedByRookOrQueen(t *testing.T) {
	board := NewBoard()

	board.Pieces[64] = wR
	board.Pieces[98] = bQ

	tests := []struct {
		square   Square
		side     Color
		expected bool
	}{
		{square: 54, side: White, expected: true},
		{square: 34, side: White, expected: true},
		{square: 24, side: White, expected: true},

		{square: 119, side: White, expected: false},
		{square: 83, side: White, expected: false},
		{square: 93, side: White, expected: false},

		{square: 78, side: Black, expected: true},
		{square: 58, side: Black, expected: true},
		{square: 94, side: Black, expected: true},
		{square: 90, side: Black, expected: false},
	}

	for _, test := range tests {
		result := board.isAttackedByRookOrQueen(test.square, test.side)

		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}

	}

}

func TestIsSquareAttackedByBishopOrQueen(t *testing.T) {
	board := NewBoard()

	board.Pieces[64] = wB
	board.Pieces[98] = bQ

	tests := []struct {
		square   Square
		side     Color
		expected bool
	}{
		{square: 73, side: White, expected: true},
		{square: 82, side: White, expected: true},
		{square: 91, side: White, expected: true},

		{square: 100, side: White, expected: false},
		{square: 52, side: White, expected: false},
		{square: 61, side: White, expected: false},

		{square: 87, side: Black, expected: true},
		{square: 76, side: Black, expected: true},
		{square: 65, side: Black, expected: true},
		{square: 10, side: Black, expected: false},
	}

	for _, test := range tests {
		result := board.isAttackedByBishopOrQueen(test.square, test.side)

		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}

	}

}

func TestIsAttackedByKing(t *testing.T) {
	board := NewBoard()

	board.Pieces[21] = bK
	board.Pieces[64] = wK

	tests := []struct {
		square   Square
		side     Color
		expected bool
	}{
		{square: 22, side: Black, expected: true},
		{square: 31, side: Black, expected: true},
		{square: 32, side: Black, expected: true},

		{square: 20, side: Black, expected: false},
		{square: 11, side: Black, expected: false},
		{square: 23, side: Black, expected: false},

		{square: 53, side: White, expected: true},
		{square: 54, side: White, expected: true},
		{square: 55, side: White, expected: true},
		{square: 63, side: White, expected: true},
		{square: 65, side: White, expected: true},
		{square: 75, side: White, expected: true},
		{square: 74, side: White, expected: true},
		{square: 73, side: White, expected: true},
		{square: 32, side: White, expected: false},
	}

	for _, test := range tests {
		result := board.isAttackedByKing(test.square, test.side)

		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}

	}

}

func TestIsAttacked(t *testing.T) {
	board := NewBoard()
	board.ParseFEN(START)

	tests := []struct {
		square   Square
		side     Color
		expected bool
	}{
		{square: 55, side: Black, expected: false},
		{square: 55, side: White, expected: false},

		// black knight
		{square: 71, side: Black, expected: true},
		{square: 73, side: Black, expected: true},

		// black pawn
		{square: 73, side: Black, expected: true},
		{square: 76, side: Black, expected: true},
	}

	for _, test := range tests {
		result := board.IsAttacked(test.square, test.side)

		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}

	}

}
