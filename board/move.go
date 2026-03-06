package board

type Move struct{
	move int
	score int
}

// 0000 0000 0000 0000 0000 0111 1111 -> From Square (7 bits)
// 0000 0000 0000 0011 1111 1000 0000 -> To Square (7 bits)
// 0000 0000 0011 1100 0000 0000 0000 -> Captured piece (4 bits)
// 0000 0000 0100 0000 0000 0000 0000 -> Enpassant move (1 bit)
// 0000 0000 1000 0000 0000 0000 0000 -> Pawn start (1 bit)
// 0000 1111 1000 0000 0000 0000 0000 -> Promotion piece (4 bit)
// 0001 0000 0000 0000 0000 0000 0000 -> Castle move (1 bit)


