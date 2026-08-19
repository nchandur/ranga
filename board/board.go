package board

import "fmt"

// represents complete chess board state
type Board struct {
	PieceBitBoards [12]BitBoard // bitboards for each piece
	Occupancies    [3]BitBoard  // combined bitboards for white, black and all pieces
	Mailbox        [64]Piece
	Side           Color  // side to move
	EnPassant      Square // en passant square
	Castle                // castle rights
	Ply            int    // current half-move ply
	FiftyMove      int    // halfmove clock for 50-move draw rule
	ZobristKey     uint64 // unique hash key for position
}

// creates a new instance of board
func NewBoard() Board {
	var res Board

	res.PieceBitBoards = [12]BitBoard{}
	res.Occupancies = [3]BitBoard{}

	for i := range 64 {
		res.Mailbox[i] = Empty
	}

	res.Side = Both
	res.EnPassant = NoSquare
	res.Castle = 0
	res.Ply = 0
	res.FiftyMove = 0
	res.ZobristKey = 0

	return res
}

func (b *Board) Clear() {
	b.PieceBitBoards = [12]BitBoard{}
	b.Occupancies = [3]BitBoard{}

	for i := range 64 {
		b.Mailbox[i] = Empty
	}
	b.Side = Both
	b.EnPassant = NoSquare
	b.Castle = 0
	b.Ply = 0
	b.FiftyMove = 0
	b.ZobristKey = 0

}

// add piece to the board
func (b *Board) AddPiece(piece Piece, sq Square) {
	b.PieceBitBoards[piece].SetBit(sq)

	if piece <= WK {
		b.Occupancies[White].SetBit(sq)
	} else if piece > WK && piece <= BK {
		b.Occupancies[Black].SetBit(sq)
	}

	b.Occupancies[Both].SetBit(sq)
	b.Mailbox[sq] = piece

}

// remove piece from the board
func (b *Board) RemovePiece(sq Square) {
	piece := b.Mailbox[sq]
	b.PieceBitBoards[piece].PopBit(sq)

	if piece <= WK {
		b.Occupancies[White].PopBit(sq)
	} else if piece > WK && piece <= BK {
		b.Occupancies[Black].PopBit(sq)
	}

	b.Occupancies[Both].PopBit(sq)
	b.Mailbox[sq] = Empty

}

func (b *Board) Print() {
	fmt.Println("     a   b   c   d   e   f   g   h")
	fmt.Println("   +---+---+---+---+---+---+---+---+")

	for rank := 0; rank <= 7; rank++ {
		fmt.Printf(" %d |", 8-rank)

		for file := 0; file <= 7; file++ {
			sq := FRtoSq(rank, file)
			piece := b.Mailbox[sq]
			fmt.Printf(" %c |", PieceChar[piece])
		}

		fmt.Printf(" %d\n", 8-rank)
		fmt.Println("   +---+---+---+---+---+---+---+---+")
	}

	fmt.Println("     a   b   c   d   e   f   g   h")

	fmt.Printf("\nSide to Move : %c\n", b.Side)
	fmt.Printf("En Passant   : %s\n", b.EnPassant)
	fmt.Printf("Castling     : %s\n", b.Castle)
	fmt.Printf("Zobrist      : 0x%x\n", b.ZobristKey)
}
