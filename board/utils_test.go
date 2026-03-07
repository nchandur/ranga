package board

import "testing"

func TestFRToSq(t *testing.T) {
	tests := []struct {
		file     File
		rank     Rank
		expected Square
	}{
		{file: A, rank: Eight, expected: A8},
		{file: H, rank: Four, expected: H4},
		{file: B, rank: Two, expected: B2},
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
		idx      int
		expected int
	}{
		{idx: 10, expected: 65},
		{idx: 46, expected: 21},
		{idx: 21, expected: 0},
		{idx: 119, expected: 65},
	}

	for _, test := range tests {
		output := Fr120To64(test.idx)

		if output != test.expected {
			t.Errorf("expected: %d, output: %d", test.expected, output)
		}

	}

}

func TestFr64To120(t *testing.T) {
	tests := []struct {
		idx      int
		expected int
	}{
		{idx: 0, expected: 21},
		{idx: 42, expected: 73},
		{idx: 51, expected: 84},
	}

	for _, test := range tests {
		output := Fr64To120(test.idx)

		if output != test.expected {
			t.Errorf("expected: %d, output: %d", test.expected, output)
		}

	}

}

func TestFr120ToFile(t *testing.T) {
	tests := []struct {
		idx      int
		expected File
	}{
		{idx: 91, expected: A},
		{idx: 37, expected: G},
		{idx: 55, expected: E},
	}

	for _, test := range tests {
		output := Fr120ToFile(test.idx)

		if output != test.expected {
			t.Errorf("expected: %d, output: %d", test.expected, output)
		}

	}

}

func TestFr120ToRank(t *testing.T) {
	tests := []struct {
		idx      int
		expected Rank
	}{
		{idx: 91, expected: Eight},
		{idx: 37, expected: Two},
		{idx: 55, expected: Four},
	}

	for _, test := range tests {
		output := Fr120ToRank(test.idx)

		if output != test.expected {
			t.Errorf("expected: %d, output: %d", test.expected, output)
		}

	}

}
