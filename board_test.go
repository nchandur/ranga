package main

import "testing"

func TestBoardReset(t *testing.T) {
	board := NewBoard()

	board.Pieces[10] = 5
	board.Material[0] = 200
	board.PieceNumber[3] = 2
	board.SideToMove = 1
	board.EnPassant = 25
	board.FiftyMove = 10
	board.Ply = 3
	board.HistoryPly = 5
	board.Castling = 7
	board.Hash = 123456

	board.Reset()

	for i := range 64 {
		sq := Fr64To120[i]
		if board.Pieces[sq] != Empty {
			t.Errorf("Pieces[%d] expected Empty, got %d", sq, board.Pieces[sq])
		}
	}

	for i := range board.PieceList {
		if board.PieceList[i] != Square(Empty) {
			t.Errorf("PieceList[%d] expected Empty, got %d", i, board.PieceList[i])
		}
	}

	for i := range board.Material {
		if board.Material[i] != 0 {
			t.Errorf("Material[%d] expected 0, got %d", i, board.Material[i])
		}
	}

	for i := range board.PieceNumber {
		if board.PieceNumber[i] != 0 {
			t.Errorf("PieceNumber[%d] expected 0, got %d", i, board.PieceNumber[i])
		}
	}

	if board.SideToMove != Both {
		t.Errorf("SideToMove expected %d, got %d", Both, board.SideToMove)
	}

	if board.EnPassant != NoSquare {
		t.Errorf("EnPassant expected %d, got %d", NoSquare, board.EnPassant)
	}

	if board.FiftyMove != 0 {
		t.Errorf("FiftyMove expected 0, got %d", board.FiftyMove)
	}

	if board.HistoryPly != 0 {
		t.Errorf("HistoryPly expected 0, got %d", board.HistoryPly)
	}

	if board.Ply != 0 {
		t.Errorf("Ply expected 0, got %d", board.Ply)
	}

	if board.Castling != 0 {
		t.Errorf("Castling expected 0, got %d", board.Castling)
	}

	if board.Hash != 0 {
		t.Errorf("Hash expected 0, got %d", board.Hash)
	}

	if board.MoveList[board.Ply] != 0 {
		t.Errorf("MoveList[%d] expected 0, got %d", board.Ply, board.MoveList[board.Ply])
	}
}

func TestUpdatePieceList(t *testing.T) {
	board := NewBoard()
	board.Reset()

	sq1 := Fr64To120[0]
	sq2 := Fr64To120[10]
	sq3 := Fr64To120[20]

	board.Pieces[sq1] = wP
	board.Pieces[sq2] = bP
	board.Pieces[sq3] = wN

	board.UpdatePieceList()

	if board.PieceNumber[wP] != 1 {
		t.Errorf("expected 1 white pawN, got %d", board.PieceNumber[wP])
	}
	if board.PieceNumber[bP] != 1 {
		t.Errorf("expected 1 black pawN, got %d", board.PieceNumber[bP])
	}
	if board.PieceNumber[wN] != 1 {
		t.Errorf("expected 1 white knight, got %d", board.PieceNumber[wN])
	}

	wPIdx := (int(wP) * 10) + (board.PieceNumber[wP] - 1)
	if board.PieceList[wPIdx] != Square(sq1) {
		t.Errorf("expected white pawN at square %d, got %d", sq1, board.PieceList[wPIdx])
	}

	bPIdx := (int(bP) * 10) + (board.PieceNumber[bP] - 1)
	if board.PieceList[bPIdx] != Square(sq2) {
		t.Errorf("expected black pawN at square %d, got %d", sq2, board.PieceList[bPIdx])
	}

	wNIdx := (int(wN) * 10) + (board.PieceNumber[wN] - 1)
	if board.PieceList[wNIdx] != Square(sq3) {
		t.Errorf("expected white knight at square %d, got %d", sq3, board.PieceList[wNIdx])
	}

	expectedWhiteMaterial := pieceValue[wP] + pieceValue[wN]
	expectedBlackMaterial := pieceValue[bP]

	if board.Material[White] != expectedWhiteMaterial {
		t.Errorf("expected white material %d, got %d",
			expectedWhiteMaterial, board.Material[White])
	}
	if board.Material[Black] != expectedBlackMaterial {
		t.Errorf("expected black material %d, got %d",
			expectedBlackMaterial, board.Material[Black])
	}
}

func TestValidateFEN(t *testing.T) {
	board := NewBoard()

	tests := []struct {
		fen       string
		shouldErr bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},
		{"8/8/8/8/8/8/8/8 b - - 0 1", false},
		{"8/8/8/8/8/8/8/8 w - 0", true},
		{"8/8/8/8/8/8/8 w KQkq - 0 1", true},
		{"8/8/8/8/8/8/8/X w KQkq - 0 1", true},
		{"8/8/8/8/8/8/8/8 x KQkq - 0 1", true},
		{"8/8/8/8/8/8/8/8 w ABC - 0 1", true},
		{"8/8/8/8/8/8/8/8 w KQkq i9 0 1", true},
		{"8/8/8/8/8/8/8/8 w KQkq - -1 1", true},
		{"8/8/8/8/8/8/8/8 w KQkq - 0 0", true},
	}

	for _, test := range tests {
		err := board.validateFEN(test.fen)
		if test.shouldErr && err == nil {
			t.Errorf("expected error for FEN: %q, but got none", test.fen)
		}
		if !test.shouldErr && err != nil {
			t.Errorf("did not expect error for FEN: %q, but got: %v", test.fen, err)
		}
	}
}
