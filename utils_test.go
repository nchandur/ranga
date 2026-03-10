package main

import "testing"

func TestFR2Square(t *testing.T) {
	tests := []struct {
		file     File
		rank     Rank
		expected Square
	}{
		{file: FileA, rank: Rank1, expected: 21},
		{file: FileE, rank: Rank5, expected: 65},
		{file: FileG, rank: Rank3, expected: 47},
	}

	for _, test := range tests {
		output := FR2Square(test.file, test.rank)

		if output != test.expected {
			t.Errorf("Expected: %d, Output: %d", test.expected, output)
		}
	}

}
