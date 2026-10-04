package search

import (
	"ranga/internal/board"
	"testing"
)

func TestMoveOrdering(t *testing.T) {
	mTT := board.NewMove(board.E2, board.E4, board.WP, board.Empty, false, false, false, false)
	mPV := board.NewMove(board.D2, board.D4, board.WP, board.Empty, false, false, false, false)
	mPromoQ := board.NewMove(board.E7, board.E8, board.WP, board.WQ, false, false, false, false)
	mKiller1 := board.NewMove(board.G1, board.F3, board.WN, board.Empty, false, false, false, false)
	mKiller2 := board.NewMove(board.B1, board.C3, board.WN, board.Empty, false, false, false, false)
	mHistory := board.NewMove(board.A2, board.A3, board.WP, board.Empty, false, false, false, false)

	t.Run("scoreMove assigns scores following priority tiers", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
		b.Ply = 0

		if score := s.scoreMove(&b, mTT, mTT); score != 30000 {
			t.Errorf("TT move score = %d; want 30000", score)
		}

		s.PV.ScorePV = true
		s.PV.Table[0][b.Ply] = mPV
		if score := s.scoreMove(&b, mPV, board.NOMOVE); score != 20000 {
			t.Errorf("PV move score = %d; want 20000", score)
		}
		if s.PV.ScorePV {
			t.Errorf("ScorePV was not reset to false after scoring PV move")
		}

		mCapture := board.NewMove(board.E4, board.D5, board.WP, board.Empty, true, false, false, false)
		expectedCapScore := MVVLVA[board.WP][board.BP] + 10000
		if score := s.scoreMove(&b, mCapture, board.NOMOVE); score != expectedCapScore {
			t.Errorf("Capture score = %d; want %d", score, expectedCapScore)
		}

		if score := s.scoreMove(&b, mPromoQ, board.NOMOVE); score != 9500 {
			t.Errorf("Queen promo score = %d; want 9500", score)
		}

		s.Killers[0][b.Ply] = mKiller1
		s.Killers[1][b.Ply] = mKiller2
		if score := s.scoreMove(&b, mKiller1, board.NOMOVE); score != 9000 {
			t.Errorf("Killer 1 score = %d; want 9000", score)
		}
		if score := s.scoreMove(&b, mKiller2, board.NOMOVE); score != 8000 {
			t.Errorf("Killer 2 score = %d; want 8000", score)
		}

		s.History[board.WP][board.A3] = 450
		if score := s.scoreMove(&b, mHistory, board.NOMOVE); score != 450 {
			t.Errorf("History score = %d; want 450", score)
		}
	})
	t.Run("scoreMove handles en passant captures", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2")

		mEP := board.NewMove(board.E5, board.D6, board.WP, board.Empty, true, false, true, false)
		score := s.scoreMove(&b, mEP, board.NOMOVE)

		wantScore := MVVLVA[board.WP][board.BP] + 10000
		if score != wantScore {
			t.Errorf("En passant capture score = %d; want %d", score, wantScore)
		}
	})
	t.Run("sortMove arranges moves in descending order of score", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
		b.Ply = 0

		s.Killers[0][b.Ply] = mKiller1
		s.History[board.WP][board.A3] = 120

		mCapture := board.NewMove(board.E4, board.D5, board.WP, board.Empty, true, false, false, false)

		var ml board.MoveList
		ml.AddMove(mHistory)
		ml.AddMove(mTT)
		ml.AddMove(mKiller1)
		ml.AddMove(mCapture)

		s.sortMove(&b, &ml, mTT)

		expected := []board.Move{mTT, mCapture, mKiller1, mHistory}
		for i, want := range expected {
			if ml.Moves[i] != want {
				t.Errorf("Move at index %d = %v; want %v", i, ml.Moves[i], want)
			}
		}
	})
}
