package search

import "math"

func initialiseLMR() {
	for d := 1; d < 64; d++ {
		for m := 1; m < 64; m++ {
			lmrTable[d][m] = int(0.75 + math.Log(float64(d))*math.Log(float64(m))/2.25)
		}
	}
}

func init() {
	initialiseLMR()
}
