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
	Network
	acc perspectives
}

func NewRandom() *NNUE {
	n := &NNUE{}
	n.Network.Randomize()
	return n
}

// rebuilds both accumulators from scratch
func (n *NNUE) Reset(b *board.Board) {
	n.acc.White = n.FeatureBias
	n.acc.Black = n.FeatureBias

	for sq := range board.Square(64) {
		piece := b.Mailbox[sq]
		if piece == board.Empty {
			continue
		}
		n.addPiece(piece, sq)
	}
}

func (n *NNUE) Preserve() Snapshot {
	return Snapshot(n.acc)
}

func (n *NNUE) Restore(s Snapshot) {
	n.acc = perspectives(s)
}

// applies incremental accumulator changes for a move that has already been confirmed as legal
func (n *NNUE) Update(before *board.Board, move board.Move) {
	source, target := move.Source(), move.Target()
	piece := move.Piece()

	n.removePiece(piece, source)

	if move.IsCapture() && !move.IsEnpass() {
		n.removePiece(before.Mailbox[target], target)
	}

	if move.IsEnpass() {
		if before.Side == board.White {
			n.removePiece(board.BP, target+8)
		} else {
			n.removePiece(board.WP, target-8)
		}
	}

	if move.Promoted() != board.Empty {
		n.addPiece(move.Promoted(), target)
	} else {
		n.addPiece(piece, target)
	}

	if move.IsCastle() {
		switch target {
		case board.G1:
			n.removePiece(board.WR, board.H1)
			n.addPiece(board.WR, board.F1)
		case board.C1:
			n.removePiece(board.WR, board.A1)
			n.addPiece(board.WR, board.D1)
		case board.G8:
			n.removePiece(board.BR, board.H8)
			n.addPiece(board.BR, board.F8)
		case board.C8:
			n.removePiece(board.BR, board.A8)
			n.addPiece(board.BR, board.D8)
		}
	}
}

func (n *NNUE) Evaluate(b *board.Board) int {
	if b.Side == board.White {
		return int(n.Network.Evaluate(&n.acc.White, &n.acc.Black))
	}
	return int(n.Network.Evaluate(&n.acc.Black, &n.acc.White))
}

func (n *NNUE) addPiece(piece board.Piece, sq board.Square) {
	isWhite := int(piece) <= 5

	n.acc.White.AddFeature(&n.Network, FeatureIndex(true, isWhite, piece, sq))
	n.acc.Black.AddFeature(&n.Network, FeatureIndex(false, !isWhite, piece, sq))
}

func (n *NNUE) removePiece(piece board.Piece, sq board.Square) {
	isWhite := int(piece) <= 5

	n.acc.White.RemoveFeature(&n.Network, FeatureIndex(true, isWhite, piece, sq))
	n.acc.Black.RemoveFeature(&n.Network, FeatureIndex(false, !isWhite, piece, sq))
}
