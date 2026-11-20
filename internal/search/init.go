package search

import "math"

func initLMRTable() {
	for d := 1; d < 64; d++ {
		for m := 1; m < 64; m++ {
			lmrTable[d][m] = int(lmrBase + math.Log(float64(d))*math.Log(float64(m))/lmrDiv)
		}
	}
}

func init() {
	initLMRTable()
}
