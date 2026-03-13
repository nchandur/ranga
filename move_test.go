package main

import "testing"

func TestMoveFromSquare(t *testing.T) {
	tests := []struct {
		move     Move
		expected Square
	}{
		{move: NewMove(53, NoSquare, Empty, Empty, MNONE), expected: 53},
		{move: NewMove(97, NoSquare, Empty, Empty, MNONE), expected: 97},
		{move: NewMove(61, NoSquare, Empty, Empty, MNONE), expected: 61},
		{move: NewMove(119, NoSquare, Empty, Empty, MNONE), expected: 119},
	}

	for _, test := range tests {
		output := test.move.FromSquare()

		if output != test.expected {
			t.Errorf("expected %v, got %v", test.expected, output)
		}

	}

}

func TestMoveToSquare(t *testing.T) {
	tests := []struct {
		move     Move
		expected Square
	}{
		{move: NewMove(NoSquare, 53, Empty, Empty, MNONE), expected: 53},
		{move: NewMove(NoSquare, 97, Empty, Empty, MNONE), expected: 97},
		{move: NewMove(NoSquare, 61, Empty, Empty, MNONE), expected: 61},
		{move: NewMove(NoSquare, 119, Empty, Empty, MNONE), expected: 119},
	}

	for _, test := range tests {
		output := test.move.ToSquare()

		if output != test.expected {
			t.Errorf("expected %v, got %v", test.expected, output)
		}

	}

}

func TestMoveCaptured(t *testing.T) {
	tests := []struct {
		move     Move
		expected Piece
	}{
		{move: NewMove(NoSquare, NoSquare, Empty, Empty, MNONE), expected: Empty},
		{move: NewMove(NoSquare, NoSquare, bP, Empty, MNONE), expected: bP},
		{move: NewMove(NoSquare, NoSquare, wQ, Empty, MNONE), expected: wQ},
		{move: NewMove(NoSquare, NoSquare, bB, Empty, MNONE), expected: bB},
		{move: NewMove(NoSquare, NoSquare, wR, Empty, MNONE), expected: wR},
	}

	for _, test := range tests {
		output := test.move.Captured()

		if output != test.expected {
			t.Errorf("expected %v, got %v", test.expected, output)
		}

	}

}

func TestMovePromoted(t *testing.T) {
	tests := []struct {
		move     Move
		expected Piece
	}{
		{move: NewMove(NoSquare, NoSquare, Empty, Empty, MNONE), expected: Empty},
		{move: NewMove(NoSquare, NoSquare, Empty, bN, MNONE), expected: bN},
		{move: NewMove(NoSquare, NoSquare, Empty, wQ, MNONE), expected: wQ},
		{move: NewMove(NoSquare, NoSquare, Empty, bB, MNONE), expected: bB},
		{move: NewMove(NoSquare, NoSquare, Empty, wR, MNONE), expected: wR},
	}

	for _, test := range tests {
		output := test.move.Promoted()

		if output != test.expected {
			t.Errorf("expected %v, got %v", test.expected, output)
		}

	}

}
