package nnue

func init() {
	for p := range 2 {
		for sq := range 64 {
			rel := sq
			if p == 0 {
				rel ^= 56
			}
			file := rel & 7
			if file >= 4 {
				file ^= 7
			}
			b := inputBucketLayout[(rel>>3)*4+file]
			if b < 0 || b >= InputBuckets {
				panic("nnue: inputBucketLayout entry out of range of InputBuckets")
			}
			bucketTable[p][sq] = uint8(b)
		}
	}
}
