package nnue

import "testing"

func Test_Utils(t *testing.T) {
	t.Run("clamp bounds int16 values correctly", func(t *testing.T) {
		tests := []struct {
			name         string
			n, low, high int16
			want         int16
		}{
			{name: "within range", n: 10, low: 0, high: 255, want: 10},
			{name: "below lower bound", n: -50, low: 0, high: 255, want: 0},
			{name: "above upper bound", n: 300, low: 0, high: 255, want: 255},
			{name: "exactly at lower bound", n: 0, low: 0, high: 255, want: 0},
			{name: "exactly at upper bound", n: 255, low: 0, high: 255, want: 255},
			{name: "negative boundaries", n: -10, low: -20, high: -5, want: -10},
			{name: "exceeding negative upper bound", n: 0, low: -20, high: -5, want: -5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := clamp(tt.n, tt.low, tt.high); got != tt.want {
					t.Errorf("clamp(%d, %d, %d) = %d; want %d",
						tt.n, tt.low, tt.high, got, tt.want)
				}
			})
		}
	})

	t.Run("screlu evaluates squared clipped ReLU accurately", func(t *testing.T) {
		tests := []struct {
			name string
			x    int16
			want int32
		}{
			{name: "negative input clamps to 0", x: -100, want: 0},
			{name: "zero input returns 0", x: 0, want: 0},
			{name: "small positive squared", x: 10, want: 100},
			{name: "mid positive squared", x: 50, want: 2500},
			{name: "boundary QA returns QA^2", x: QA, want: int32(QA) * int32(QA)},
			{name: "saturated above QA returns QA^2", x: QA + 50, want: int32(QA) * int32(QA)},
			{name: "large input saturates at QA^2", x: 30000, want: int32(QA) * int32(QA)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := screlu(tt.x); got != tt.want {
					t.Errorf("screlu(%d) = %d; want %d", tt.x, got, tt.want)
				}
			})
		}
	})
}
