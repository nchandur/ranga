package main

import "testing"

func TestAddCaptureMove(t *testing.T) {
	tests := []struct {
		ply           int
		startIndex    int
		moveToAdd     int
		expectedIndex int
		expectedMove  int
	}{
		{
			ply:           0,
			startIndex:    0,
			moveToAdd:     42,
			expectedIndex: 1,
			expectedMove:  42,
		},
		{
			ply:           5,
			startIndex:    3,
			moveToAdd:     99,
			expectedIndex: 4,
			expectedMove:  99,
		},
		{
			ply:           2,
			startIndex:    7,
			moveToAdd:     55,
			expectedIndex: 8,
			expectedMove:  55,
		},
	}

	for _, test := range tests {
		board := NewBoard()
		board.Ply = test.ply
		board.MoveListStart[test.ply+1] = test.startIndex

		board.addCaptureMove(Move{move: test.moveToAdd})

		// Check the move is added
		if board.MoveList[test.startIndex] != test.expectedMove {
			t.Errorf("expected: %d, output: %d", test.expectedMove, board.MoveList[test.startIndex])
		}

		// Check MoveListStart incremented
		if board.MoveListStart[test.ply+1] != test.expectedIndex {
			t.Errorf("expected: %d, output: %d", test.expectedIndex, board.MoveListStart[test.ply+1])
		}

	}
}

func TestAddQuietMove(t *testing.T) {
	tests := []struct {
		ply           int
		startIndex    int
		moveToAdd     int
		expectedIndex int
		expectedMove  int
	}{
		{
			ply:           0,
			startIndex:    0,
			moveToAdd:     12,
			expectedIndex: 1,
			expectedMove:  12,
		},
		{
			ply:           3,
			startIndex:    5,
			moveToAdd:     77,
			expectedIndex: 6,
			expectedMove:  77,
		},
		{
			ply:           2,
			startIndex:    10,
			moveToAdd:     99,
			expectedIndex: 11,
			expectedMove:  99,
		},
	}

	for _, test := range tests {
		board := NewBoard()
		board.Ply = test.ply
		board.MoveListStart[test.ply+1] = test.startIndex

		board.addQuietMove(Move{move: test.moveToAdd})

		// Check the move is added
		if board.MoveList[test.startIndex] != test.expectedMove {
			t.Errorf("expected: %d, output: %d", test.expectedMove, board.MoveList[test.startIndex])
		}

		// Check MoveListStart incremented
		if board.MoveListStart[test.ply+1] != test.expectedIndex {
			t.Errorf("expected: %d, output: %d", test.expectedIndex, board.MoveListStart[test.ply+1])
		}

	}
}

func TestAddEnpassantMove(t *testing.T) {
	tests := []struct {
		ply           int
		startIndex    int
		moveToAdd     int
		expectedIndex int
		expectedMove  int
	}{
		{
			ply:           0,
			startIndex:    0,
			moveToAdd:     22,
			expectedIndex: 1,
			expectedMove:  22,
		},
		{
			ply:           4,
			startIndex:    2,
			moveToAdd:     55,
			expectedIndex: 3,
			expectedMove:  55,
		},
		{
			ply:           2,
			startIndex:    7,
			moveToAdd:     99,
			expectedIndex: 8,
			expectedMove:  99,
		},
	}

	for _, test := range tests {
		board := NewBoard()
		board.Ply = test.ply
		board.MoveListStart[test.ply+1] = test.startIndex

		board.addQuietMove(Move{move: test.moveToAdd})

		// Check the move is added
		if board.MoveList[test.startIndex] != test.expectedMove {
			t.Errorf("expected: %d, output: %d", test.expectedMove, board.MoveList[test.startIndex])
		}

		// Check MoveListStart incremented
		if board.MoveListStart[test.ply+1] != test.expectedIndex {
			t.Errorf("expected: %d, output: %d", test.expectedIndex, board.MoveListStart[test.ply+1])
		}

	}

}
