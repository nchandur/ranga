package nnue

import (
	"ranga/internal/board"
)

type perspectives struct {
	White, Black Accumulator
}

// opaque copy of accumulator state
type Snapshot perspectives

type NNUE struct {
	*Network
	acc perspectives
}

func NewRandom() *NNUE {
	n := &NNUE{Network: &Network{}}
	n.Network.Randomize()
	return n
}

// rebuilds both accumulators from scratch
func (n *NNUE) Reset(b *board.Board) {
	n.acc.White = n.FeatureBias
	n.acc.Black = n.FeatureBias

	wKSq := board.Square(b.PieceBitBoards[board.WK].GetLSB())
	bKSq := board.Square(b.PieceBitBoards[board.BK].GetLSB())

	for sq := range board.Square(64) {
		piece := b.Mailbox[sq]
		if piece == board.Empty {
			continue
		}
		n.addPiece(piece, sq, wKSq, bKSq)
	}
}

func (n *NNUE) Preserve() Snapshot {
	return Snapshot(n.acc)
}

func (n *NNUE) Restore(s Snapshot) {
	n.acc = perspectives(s)
}

// applies incremental accumulator changes using fused loops
func (n *NNUE) Update(before *board.Board, move board.Move) {
	source, target := move.Source(), move.Target()
	movingPiece := before.Mailbox[source]

	wKSq := board.Square(before.PieceBitBoards[board.WK].GetLSB())
	bKSq := board.Square(before.PieceBitBoards[board.BK].GetLSB())

	newWKSq, newBKSq := wKSq, bKSq
	if movingPiece == board.WK {
		newWKSq = target
	} else if movingPiece == board.BK {
		newBKSq = target
	}

	// castling
	if move.IsCastle() {
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
		wSub1, bSub1 := getIndices(movingPiece, source, wKSq, bKSq)
		wSub2, bSub2 := getIndices(rPiece, rSource, wKSq, bKSq)

		n.acc.White.Add2Sub2(n.Network, wAdd1, wAdd2, wSub1, wSub2)
		n.acc.Black.Add2Sub2(n.Network, bAdd1, bAdd2, bSub1, bSub2)
		return
	}

	addPiece := movingPiece
	if move.IsPromotion() {
		addPiece = move.Promoted(before.Side)
	}

	// captures and enpassant
	if move.IsCapture() || move.IsEnpass() {
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
		wSub1, bSub1 := getIndices(movingPiece, source, wKSq, bKSq)
		wSub2, bSub2 := getIndices(capPiece, capTarget, wKSq, bKSq)

		n.acc.White.AddSub2(n.Network, wAdd, wSub1, wSub2)
		n.acc.Black.AddSub2(n.Network, bAdd, bSub1, bSub2)
		return
	}

	// quiet moves
	wAdd, bAdd := getIndices(addPiece, target, newWKSq, newBKSq)
	wSub, bSub := getIndices(movingPiece, source, wKSq, bKSq)
	n.acc.White.AddSub(n.Network, wAdd, wSub)
	n.acc.Black.AddSub(n.Network, bAdd, bSub)
}

func (n *NNUE) Evaluate(b *board.Board) int {
	bucket := outputBucket(b)
	if b.Side == board.White {
		return int(n.Network.Evaluate(&n.acc.White, &n.acc.Black, bucket))
	}
	return int(n.Network.Evaluate(&n.acc.Black, &n.acc.White, bucket))
}

func (n *NNUE) addPiece(piece board.Piece, sq board.Square, wKSq, bKSq board.Square) {
	isWhite := int(piece) <= 5
	wIdx := FeatureIndex(true, isWhite, piece, sq, wKSq)
	bIdx := FeatureIndex(false, !isWhite, piece, sq, bKSq)
	n.acc.White.AddFeature(n.Network, wIdx)
	n.acc.Black.AddFeature(n.Network, bIdx)
}
