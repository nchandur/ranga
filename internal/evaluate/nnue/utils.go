package nnue

func screlu(x int16) int32 {
	y := int32(clamp(int(x), 0, QA))
	return y * y
}

// helper function to clamp bonus between MAX_HISTORY
func clamp(n, low, high int) int {
	if n <= low {
		return low
	}

	if n >= high {
		return high
	}

	return n
}
