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

		{square: 71, side: Black, expected: true},
		{square: 73, side: Black, expected: true},
		{square: 84, side: Black, expected: true},
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
