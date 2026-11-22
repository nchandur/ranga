package nnue

import "ranga/internal/board"

const (
	deltaNone    uint8 = iota // null move
	deltaQuiet                // 1 add, 1 sub
	deltaCapture              // 1 add, 2 sub
	deltaCastle               // 2 add, 2 sub
)

type entry struct {
	acc  [2]Accumulator
	kind uint8
	add  [2][2]uint16 // [perspective][i] feature indices
	sub  [2][2]uint16
	ok   [2]bool // acc[p] is up to date
}

type NNUE struct {
	*Network
	stack [StackSize]entry
	ply   int
}

func NewRandom() *NNUE {
	n := &NNUE{Network: &Network{}}
	n.Network.Randomize()
	return n
}

// rebuilds root entry from scratch and resets stack
func (n *NNUE) Reset(b *board.Board) {
	n.ply = 0
	e := &n.stack[0]
	e.kind = deltaNone
	e.ok = [2]bool{true, true}
	e.acc[int(board.White)] = n.FeatureBias
	e.acc[int(board.Black)] = n.FeatureBias

	wKSq := board.Square(b.PieceBitBoards[board.WK].GetLSB())
	bKSq := board.Square(b.PieceBitBoards[board.BK].GetLSB())

	for sq := range board.Square(64) {
		piece := b.Mailbox[sq]
		if piece == board.Empty {
			continue
		}
		isWhite := int(piece) <= 5
		e.acc[int(board.White)].AddFeature(n.Network, FeatureIndex(true, isWhite, piece, sq, wKSq))
		e.acc[int(board.Black)].AddFeature(n.Network, FeatureIndex(false, !isWhite, piece, sq, bKSq))
	}
}

// records the move's delta in a new stack entry
func (n *NNUE) Update(before *board.Board, move board.Move) {
	source, target := move.Source(), move.Target()
	movingPiece := before.Mailbox[source]

	wKSq := board.Square(before.PieceBitBoards[board.WK].GetLSB())
	bKSq := board.Square(before.PieceBitBoards[board.BK].GetLSB())
	newWKSq, newBKSq := wKSq, bKSq

	switch movingPiece {
	case board.WK:
		newWKSq = target
	case board.BK:
		newBKSq = target
	}

	n.ply++
	e := &n.stack[n.ply]
	e.ok = [2]bool{}

	subWKSq, subBKSq := wKSq, bKSq
	refreshW, refreshB := false, false
	if movingPiece == board.WK && needsRefresh(true, wKSq, newWKSq) {
		n.refresh(&e.acc[int(board.White)], true, before, newWKSq)
		subWKSq, refreshW = newWKSq, true
	} else if movingPiece == board.BK && needsRefresh(false, bKSq, newBKSq) {
		n.refresh(&e.acc[int(board.Black)], false, before, newBKSq)
		subBKSq, refreshB = newBKSq, true
	}

	addPiece := movingPiece
	if move.IsPromotion() {
		addPiece = move.Promoted(before.Side)
	}

	switch {
	case move.IsCastle():
		var rSource, rTarget board.Square
		var rPiece board.Piece
		switch target {
		case board.G1:
			rSource, rTarget, rPiece = board.H1, board.F1, board.WR
		case board.C1:
			rSource, rTarget, rPiece = board.A1, board.D1, board.WR
		case board.G8:
			rSource, rTarget, rPiece = board.H8, board.F8, board.BR
		case board.C8:
			rSource, rTarget, rPiece = board.A8, board.D8, board.BR
		}
		wAdd1, bAdd1 := getIndices(movingPiece, target, newWKSq, newBKSq)
		wAdd2, bAdd2 := getIndices(rPiece, rTarget, newWKSq, newBKSq)
		wSub1, bSub1 := getIndices(movingPiece, source, subWKSq, subBKSq)
		wSub2, bSub2 := getIndices(rPiece, rSource, subWKSq, subBKSq)

		e.kind = deltaCastle
		e.add = [2][2]uint16{{uint16(wAdd1), uint16(wAdd2)}, {uint16(bAdd1), uint16(bAdd2)}}
		e.sub = [2][2]uint16{{uint16(wSub1), uint16(wSub2)}, {uint16(bSub1), uint16(bSub2)}}

	case move.IsCapture() || move.IsEnpass():
		var capPiece board.Piece
		var capTarget board.Square
		if move.IsEnpass() {
			if before.Side == board.White {
				capPiece, capTarget = board.BP, target+8
			} else {
				capPiece, capTarget = board.WP, target-8
			}
		} else {
			capPiece, capTarget = before.Mailbox[target], target
		}
		wAdd, bAdd := getIndices(addPiece, target, newWKSq, newBKSq)
		wSub1, bSub1 := getIndices(movingPiece, source, subWKSq, subBKSq)
		wSub2, bSub2 := getIndices(capPiece, capTarget, subWKSq, subBKSq)

		e.kind = deltaCapture
		e.add = [2][2]uint16{{uint16(wAdd)}, {uint16(bAdd)}}
		e.sub = [2][2]uint16{{uint16(wSub1), uint16(wSub2)}, {uint16(bSub1), uint16(bSub2)}}

	default:
		wAdd, bAdd := getIndices(addPiece, target, newWKSq, newBKSq)
		wSub, bSub := getIndices(movingPiece, source, subWKSq, subBKSq)

		e.kind = deltaQuiet
		e.add = [2][2]uint16{{uint16(wAdd)}, {uint16(bAdd)}}
		e.sub = [2][2]uint16{{uint16(wSub)}, {uint16(bSub)}}
	}

	if refreshW {
		n.apply(e, int(board.White), &e.acc[int(board.White)])
		e.ok[int(board.White)] = true
	}
	if refreshB {
		n.apply(e, int(board.Black), &e.acc[int(board.Black)])
		e.ok[int(board.Black)] = true
	}
}

// null move
func (n *NNUE) PushNull() {
	n.ply++
	e := &n.stack[n.ply]
	e.kind = deltaNone
	e.ok = [2]bool{}
}

// undoes last Update or PushNull
func (n *NNUE) Pop() { n.ply-- }

func (n *NNUE) Evaluate(b *board.Board) int {
	n.materialize(int(board.White))
	n.materialize(int(board.Black))
	cur := &n.stack[n.ply]
	bucket := outputBucket(b)
	if b.Side == board.White {
		return int(n.Network.Evaluate(&cur.acc[int(board.White)], &cur.acc[int(board.Black)], bucket))
	}
	return int(n.Network.Evaluate(&cur.acc[int(board.Black)], &cur.acc[int(board.White)], bucket))
}

// brings stack[ply].acc[p] up to date
func (n *NNUE) materialize(p int) {
	base := n.ply
	for base > 0 && !n.stack[base].ok[p] {
		base--
	}
	for i := base + 1; i <= n.ply; i++ {
		e := &n.stack[i]
		n.apply(e, p, &n.stack[i-1].acc[p])
		e.ok[p] = true
	}
}

func (n *NNUE) apply(e *entry, p int, src *Accumulator) {
	dst := &e.acc[p]
	a, s := &e.add[p], &e.sub[p]
	switch e.kind {
	case deltaQuiet:
		dst.AddSubFrom(src, n.Network, int(a[0]), int(s[0]))
	case deltaCapture:
		dst.AddSub2From(src, n.Network, int(a[0]), int(s[0]), int(s[1]))
	case deltaCastle:
		dst.Add2Sub2From(src, n.Network, int(a[0]), int(a[1]), int(s[0]), int(s[1]))
	default:
		*dst = *src
	}
}

// rebuilds one perspective from the given board, indexing every piece with kSq
func (n *NNUE) refresh(acc *Accumulator, white bool, b *board.Board, kSq board.Square) {
	*acc = n.FeatureBias
	for sq := range board.Square(64) {
		p := b.Mailbox[sq]
		if p == board.Empty {
			continue
		}
		isUs := (int(p) <= 5) == white
		acc.AddFeature(n.Network, FeatureIndex(white, isUs, p, sq, kSq))
	}
}
