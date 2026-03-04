package board

var (
	SetMask   []uint64
	ClearMask []uint64
)

func InitBitMasks() {
	SetMask = make([]uint64, 64)
	ClearMask = make([]uint64, 64)

	for i := range 64 {
		SetMask[i] |= (uint64(1) << i)
		ClearMask[i] = ^SetMask[i]
	}

}
