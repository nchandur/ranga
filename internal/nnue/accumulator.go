package nnue

import "ranga/internal/board"

const (
	HiddenSize      = board.NNUEHiddenSize   // 256 per perspective
	NumPieceTypes   = 5                      // P N B R Q
	FeaturesPerKing = 64 * NumPieceTypes * 2 // 640
	InputSize       = 64 * FeaturesPerKing   // 40960
)

// wires package's update logic into board's hook variables
func init() {
	board.NNUEOnAdd = onAdd
	board.NNUEOnRemove = onRemove
	board.NNUEOnMoveDone = onMoveDone
}

// maps a piece to its HalfKP piece-type index (0..4 for P N B R Q) and color
func classify(piece board.Piece) (ptIndex int, color board.Color, isPiece bool) {
	switch piece {
	case board.WP:
		return 0, board.White, true
	case board.WN:
		return 1, board.White, true
	case board.WB:
		return 2, board.White, true
	case board.WR:
		return 3, board.White, true
	case board.WQ:
		return 4, board.White, true
	case board.WK:
		return -1, board.White, false
	case board.BP:
		return 0, board.Black, true
	case board.BN:
		return 1, board.Black, true
	case board.BB:
		return 2, board.Black, true
	case board.BR:
		return 3, board.Black, true
	case board.BQ:
		return 4, board.Black, true
	case board.BK:
		return -1, board.Black, false
	default:
		return -1, board.White, false
	}
}

// orient flips a square vertically for Black's perspective
func orient(sq board.Square, perspective board.Color) board.Square {
	if perspective == board.White {
		return sq
	}
	return sq ^ 56
}

// featureIndex computes the HalfKP feature index for one perspective
func featureIndex(kingSq, pieceSq board.Square, ptIndex int, pieceColor, perspective board.Color) int {
	ks := int(orient(kingSq, perspective))
	ps := int(orient(pieceSq, perspective))
	sameColor := 0
	if pieceColor != perspective {
		sameColor = 1
	}
	return ks*FeaturesPerKing + (ptIndex*2+sameColor)*64 + ps
}

func kingSquare(b *board.Board, perspective board.Color) board.Square {
	if perspective == board.White {
		return board.Square(b.PieceBitBoards[board.WK].GetLSB())
	}
	return board.Square(b.PieceBitBoards[board.BK].GetLSB())
}

// invoked by board.AddPiece right after a piece is placed.
func onAdd(b *board.Board, piece board.Piece, sq board.Square) {
	ptIndex, color, isPiece := classify(piece)
	if !isPiece {
		switch piece {
		case board.WK:
			b.NNUEDirty[0] = true
		case board.BK:
			b.NNUEDirty[1] = true
		}
		return
	}
	for _, perspective := range [2]board.Color{board.White, board.Black} {
		p := int(perspective)
		if b.NNUEDirty[p] {
			continue
		}
		idx := featureIndex(kingSquare(b, perspective), sq, ptIndex, color, perspective)
		addFeature(&b.NNUEAcc[p], idx)
	}
}

// invoked by board.RemovePiece right after a piece is taken off.
func onRemove(b *board.Board, piece board.Piece, sq board.Square) {
	ptIndex, color, isPiece := classify(piece)
	if !isPiece {
		switch piece {
		case board.WK:
			b.NNUEDirty[0] = true
		case board.BK:
			b.NNUEDirty[1] = true
		}
		return
	}
	for _, perspective := range [2]board.Color{board.White, board.Black} {
		p := int(perspective)
		if b.NNUEDirty[p] {
			continue
		}
		idx := featureIndex(kingSquare(b, perspective), sq, ptIndex, color, perspective)
		removeFeature(&b.NNUEAcc[p], idx)
	}
}

// runs once at the end of every successful MakeMove and repairs any perspective whose king moved this ply
func onMoveDone(b *board.Board) {
	if b.NNUEDirty[0] {
		Refresh(b, board.White)
		b.NNUEDirty[0] = false
	}
	if b.NNUEDirty[1] {
		Refresh(b, board.Black)
		b.NNUEDirty[1] = false
	}
}

// fully recomputes one perspective's accumulator from the current board state (bias + every active feature)
func Refresh(b *board.Board, perspective board.Color) {
	acc := &b.NNUEAcc[int(perspective)]
	copy(acc[:], Net.FTBias)

	ks := kingSquare(b, perspective)
	for sq := range 64 {
		piece := b.Mailbox[sq]
		if piece == board.Empty {
			continue
		}
		ptIndex, color, isPiece := classify(piece)
		if !isPiece {
			continue
		}
		idx := featureIndex(ks, board.Square(sq), ptIndex, color, perspective)
		addFeature(acc, idx)
	}
}

// recomputes both perspectives from scratch
func RefreshAll(b *board.Board) {
	Refresh(b, board.White)
	Refresh(b, board.Black)
	b.NNUEDirty[0] = false
	b.NNUEDirty[1] = false
}
