package square

import "testing"

func TestGetColor(t *testing.T) {

	tests := []struct {
		name     string
		square   Square
		expected uint8
	}{
		{name: "A1", square: A1, expected: 1},
		{name: "A2", square: B1, expected: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := GetColor(test.square)

			if output != test.expected {
				t.Errorf("expected: %d output: %d", test.expected, output)
			}

		})
	}

}
