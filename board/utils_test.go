package board

import "testing"

func TestFRToSq(t *testing.T) {

	tests := []struct {
		file     File
		rank     Rank
		expected Square
	}{
		{file: A, rank: One, expected: 21},
		{file: E, rank: Four, expected: 55},
		{file: D, rank: Eight, expected: 94},
		{file: G, rank: Six, expected: 77},
	}

	for _, test := range tests {
		output := FRToSq(test.file, test.rank)
		if output != test.expected {
			t.Errorf("expected: %d, output: %d", test.expected, output)
		}

	}

}

func TestFr120To64(t *testing.T) {
	tests := []struct {
		idx      uint8
		expected uint8
	}{
		{idx: 21, expected: 0},
		{idx: 98, expected: 63},
		{idx: 45, expected: 20},
		{idx: 119, expected: 65},
	}

	for _, test := range tests {
		output := Fr120To64(test.idx)

		if output != test.expected {
			t.Errorf("expected: %d, output: %d", test.expected, output)
		}
	}

}
