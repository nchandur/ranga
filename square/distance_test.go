package square

import "testing"

// chebyshevDist
// 	3, 3, 3, 3, 3, 3, 3, 3,
// 	3, 2, 2, 2, 2, 2, 2, 3,
// 	3, 2, 1, 1, 1, 1, 2, 3,
// 	3, 2, 1, 0, 0, 1, 2, 3,
// 	3, 2, 1, 0, 0, 1, 2, 3,
// 	3, 2, 1, 1, 1, 1, 2, 3,
// 	3, 2, 2, 2, 2, 2, 2, 3,
// 	3, 3, 3, 3, 3, 3, 3, 3,

func TestGetChebyshevDistance(t *testing.T) {

	tests := []struct {
		name     string
		square   Square
		expected int8
	}{
		{name: "A1", square: A1, expected: 3},
		{name: "B2", square: B2, expected: 2},
		{name: "C3", square: C3, expected: 1},
		{name: "D4", square: D4, expected: 0},
		{name: "E5", square: E5, expected: 0},
		{name: "F6", square: F6, expected: 1},
		{name: "G7", square: G7, expected: 2},
		{name: "H8", square: H8, expected: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := GetChebyshevDist(test.square)

			if output != test.expected {
				t.Errorf("expected %d output: %d", test.expected, output)
			}

		})
	}

}

// manhattanDist
// 	6, 5, 4, 3, 3, 4, 5, 6,
// 	5, 4, 3, 2, 2, 3, 4, 5,
// 	4, 3, 2, 1, 1, 2, 3, 4,
// 	3, 2, 1, 0, 0, 1, 2, 3,
// 	3, 2, 1, 0, 0, 1, 2, 3,
// 	4, 3, 2, 1, 1, 2, 3, 4,
// 	5, 4, 3, 2, 2, 3, 4, 5,
// 	6, 5, 4, 3, 3, 4, 5, 6,

func TestGetManhattanDistance(t *testing.T) {

	tests := []struct {
		name     string
		square   Square
		expected int8
	}{
		{name: "A1", square: A1, expected: 6},
		{name: "A2", square: A2, expected: 5},
		{name: "B3", square: B3, expected: 3},
		{name: "B4", square: B4, expected: 2},
		{name: "C4", square: C4, expected: 1},
		{name: "D4", square: D4, expected: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := GetManhattanDist(test.square)

			if output != test.expected {
				t.Errorf("expected %d output: %d", test.expected, output)
			}

		})
	}

}
