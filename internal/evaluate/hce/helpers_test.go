package hce

import (
	"ranga/internal/board"
	"testing"
)

func TestHCE_Helpers(t *testing.T) {
	h := &HCE{}

	t.Run("getPhase calculates phase score and phases properly", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN(board.START)
		phase, score := h.getPhase(&b)
		if phase != Opening {
			t.Errorf("expected Opening phase for START position, got %v", phase)
		}
		if score < PhaseScore[Opening] {
			t.Errorf("expected phaseScore >= %d, got %d", PhaseScore[Opening], score)
		}

		b.ParseFEN("4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1")
		phase, score = h.getPhase(&b)
		if phase != EndGame {
			t.Errorf("expected EndGame phase for pawn endgame, got %v", phase)
		}
		if score != 0 {
			t.Errorf("expected phaseScore 0, got %d", score)
		}

		b.ParseFEN("4k3/4r3/8/8/8/8/4R3/4K3 w - - 0 1")
		phase, score = h.getPhase(&b)
		if phase != MiddleGame {
			t.Errorf("expected MiddleGame phase, got %v (score: %d)", phase, score)
		}
	})
	t.Run("taper interpolates between opening and endgame", func(t *testing.T) {
		open := 100
		end := 20

		if got := h.taper(open, end, Opening, PhaseScore[Opening]); got != open {
			t.Errorf("taper(Opening) = %d; want %d", got, open)
		}

		if got := h.taper(open, end, EndGame, PhaseScore[EndGame]); got != end {
			t.Errorf("taper(EndGame) = %d; want %d", got, end)
		}

		midScore := PhaseScore[Opening] / 2
		expectedMid := (open + end) / 2
		if got := h.taper(open, end, MiddleGame, midScore); got != expectedMid {
			t.Errorf("taper(MiddleGame midpoint) = %d; want %d", got, expectedMid)
		}
	})
	t.Run("evalPawns scores penalties and bonuses", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("4k3/8/8/8/4P3/3P1P2/8/4K3 w - - 0 1")
		baseOpen, baseEnd := h.evalPawns(&b, board.E4, board.White)

		bDoubled := board.NewBoard()
		bDoubled.ParseFEN("4k3/8/8/8/4P3/4P3/8/4K3 w - - 0 1")
		doubledOpen, doubledEnd := h.evalPawns(&bDoubled, board.E4, board.White)

		if doubledOpen >= baseOpen {
			t.Errorf("expected doubled/isolated pawn to have lower opening score: doubled=%d, normal=%d", doubledOpen, baseOpen)
		}
		if doubledEnd >= baseEnd {
			t.Errorf("expected doubled/isolated pawn to have lower endgame score: doubled=%d, normal=%d", doubledEnd, baseEnd)
		}
	})
	t.Run("evalKnights evaluates piece square tables with rank flipping", func(t *testing.T) {
		whiteOpen, whiteEnd := h.evalKnights(board.D4, board.White)
		blackOpen, blackEnd := h.evalKnights(board.D5, board.Black)

		if whiteOpen != blackOpen || whiteEnd != blackEnd {
			t.Errorf("expected mirrored knights to evaluate symmetrically: white=(%d, %d), black=(%d, %d)",
				whiteOpen, whiteEnd, blackOpen, blackEnd)
		}
	})
	t.Run("evalBishops and evalQueens include mobility bonuses", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN("8/8/8/3q4/3B4/8/8/8 w - - 0 1")

		openB, endB := h.evalBishops(&b, board.D4, board.White)
		bishopAttacks := board.GetBishopAttacks(board.D4, b.Occupancies[board.Both]).CountBits()

		wantOpenB := PositionalScores[Opening][2][board.D4] + bishopAttacks*bishopMobility[Opening]
		wantEndB := PositionalScores[EndGame][2][board.D4] + bishopAttacks*bishopMobility[EndGame]

		if openB != wantOpenB || endB != wantEndB {
			t.Errorf("bishop eval mobility mismatch: got (%d, %d), want (%d, %d)", openB, endB, wantOpenB, wantEndB)
		}

		openQ, endQ := h.evalQueens(&b, board.D5, board.Black)
		sqFlipped := board.D5 ^ 56
		queenAttacks := board.GetQueenAttacks(board.D5, b.Occupancies[board.Both]).CountBits()

		wantOpenQ := PositionalScores[Opening][4][sqFlipped] + queenAttacks*queenMobility[Opening]
		wantEndQ := PositionalScores[EndGame][4][sqFlipped] + queenAttacks*queenMobility[EndGame]

		if openQ != wantOpenQ || endQ != wantEndQ {
			t.Errorf("queen eval mobility mismatch: got (%d, %d), want (%d, %d)", openQ, endQ, wantOpenQ, wantEndQ)
		}
	})
	t.Run("evalRooks checks file bonus logic", func(t *testing.T) {
		b := board.NewBoard()
		b.ParseFEN(board.START)

		openA1, endA1 := h.evalRooks(&b, board.A1, board.White)

		bEmpty := board.NewBoard()
		bEmpty.ParseFEN("8/8/8/8/8/8/8/R3K3 w - - 0 1")
		openEmpty, endEmpty := h.evalRooks(&bEmpty, board.A1, board.White)

		if FileMasks[board.A1] != 0 {
			if openEmpty <= openA1 || endEmpty <= endA1 {
				t.Errorf("open file rook (%d, %d) should score higher than initial closed file (%d, %d)",
					openEmpty, endEmpty, openA1, endA1)
			}
		}
	})
	t.Run("evalKings checks shield bonus logic", func(t *testing.T) {
		bSheltered := board.NewBoard()
		bSheltered.ParseFEN("4k3/8/8/8/8/8/5PPP/6K1 w - - 0 1")
		shelteredOpen, _ := h.evalKings(&bSheltered, board.G1, board.White)

		bExposed := board.NewBoard()
		bExposed.ParseFEN("4k3/8/8/8/8/8/8/6K1 w - - 0 1")
		exposedOpen, _ := h.evalKings(&bExposed, board.G1, board.White)

		if FileMasks[board.G1] != 0 {
			if shelteredOpen <= exposedOpen {
				t.Errorf("sheltered king (%d) should score higher than exposed king (%d)",
					shelteredOpen, exposedOpen)
			}
		}
	})
}
