package main

import "testing"

func TestIsPawn(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wP, expected: true},
		{piece: bP, expected: true},
		{piece: bQ, expected: false},
		{piece: wB, expected: false},
	}

	for _, test := range tests {
		output := isPawn[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsKnight(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wN, expected: true},
		{piece: bN, expected: true},
		{piece: bQ, expected: false},
		{piece: wB, expected: false},
	}

	for _, test := range tests {
		output := isKnight[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsBishop(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: bB, expected: true},
		{piece: wB, expected: true},
		{piece: bP, expected: false},
		{piece: bQ, expected: false},
	}

	for _, test := range tests {
		output := isBishop[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsRook(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wR, expected: true},
		{piece: bR, expected: true},
		{piece: bB, expected: false},
		{piece: wB, expected: false},
	}

	for _, test := range tests {
		output := isRook[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsQueen(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wQ, expected: true},
		{piece: bQ, expected: true},
		{piece: bP, expected: false},
		{piece: wB, expected: false},
	}

	for _, test := range tests {
		output := isQueen[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsKing(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wK, expected: true},
		{piece: bK, expected: true},
		{piece: bQ, expected: false},
		{piece: wB, expected: false},
	}

	for _, test := range tests {
		output := isKing[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsSlider(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wB, expected: true},
		{piece: bQ, expected: true},
		{piece: wR, expected: true},
		{piece: wP, expected: false},
		{piece: bK, expected: false},
		{piece: wN, expected: false},
	}

	for _, test := range tests {
		output := isSlider[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsMajor(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wR, expected: true},
		{piece: bQ, expected: true},
		{piece: wP, expected: false},
		{piece: bN, expected: false},
		{piece: wB, expected: false},
	}

	for _, test := range tests {
		output := isMajor[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestIsMinor(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected bool
	}{
		{piece: wR, expected: false},
		{piece: bQ, expected: false},
		{piece: wB, expected: true},
		{piece: bN, expected: true},
		{piece: wB, expected: true},
	}

	for _, test := range tests {
		output := isMinor[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %t, Output: %t", test.expected, output)
		}

	}

}

func TestPieceValue(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected int
	}{
		{piece: wR, expected: 500},
		{piece: bQ, expected: 1000},
		{piece: wP, expected: 100},
		{piece: bN, expected: 300},
		{piece: wB, expected: 300},
		{piece: Empty, expected: 0},
	}

	for _, test := range tests {
		output := pieceValue[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %d, Output: %d", test.expected, output)
		}

	}

}

func TestPieceColor(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected Color
	}{
		{piece: wR, expected: White},
		{piece: bQ, expected: Black},
		{piece: wP, expected: White},
		{piece: bN, expected: Black},
		{piece: wB, expected: White},
		{piece: Empty, expected: Both},
	}

	for _, test := range tests {
		output := pieceColor[test.piece]

		if output != test.expected {
			t.Errorf("Expected: %d, Output: %d", test.expected, output)
		}

	}

}
