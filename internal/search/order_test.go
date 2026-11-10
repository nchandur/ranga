package search

import (
	"ranga/internal/board"
	"testing"
)

func TestMoveOrdering(t *testing.T) {
	mTT := board.NewMove(board.E2, board.E4, board.MoveQuiet)
	mPV := board.NewMove(board.D2, board.D4, board.MoveQuiet)
	mPromoQ := board.NewMove(board.E7, board.E8, board.MovePromoQ)
	mKiller1 := board.NewMove(board.G1, board.F3, board.MoveQuiet)
	mKiller2 := board.NewMove(board.B1, board.C3, board.MoveQuiet)
	mHistory := board.NewMove(board.A2, board.A3, board.MoveQuiet)

	t.Run("scoreMove assigns scores following priority tiers", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
		b.Ply = 0

		if score := s.scoreMove(&b, mTT, mTT); score != SCORE_TT {
			t.Errorf("TT move score = %d; want %d", score, SCORE_TT)
		}

		s.PV.ScorePV = true
		s.PV.Table[0][b.Ply] = mPV
		if score := s.scoreMove(&b, mPV, board.NOMOVE); score != SCORE_PV {
			t.Errorf("PV move score = %d; want %d", score, SCORE_PV)
		}
		if s.PV.ScorePV {
			t.Errorf("ScorePV was not reset to false after scoring PV move")
		}

		mCapture := board.NewMove(board.E4, board.D5, board.MoveCapture)
		expectedCapScore := MVVLVA[board.WP][board.BP] + SCORE_CAPTURE
		if score := s.scoreMove(&b, mCapture, board.NOMOVE); score != expectedCapScore {
			t.Errorf("Capture score = %d; want %d", score, expectedCapScore)
		}

		if score := s.scoreMove(&b, mPromoQ, board.NOMOVE); score != SCORE_PROMOTION {
			t.Errorf("Queen promo score = %d; want %d", score, SCORE_PROMOTION)
		}

		s.stack[b.Ply].killers[0] = mKiller1
		s.stack[b.Ply].killers[1] = mKiller2
		if score := s.scoreMove(&b, mKiller1, board.NOMOVE); score != SCORE_KILLER1 {
			t.Errorf("Killer 1 score = %d; want %d", score, SCORE_KILLER1)
		}
		if score := s.scoreMove(&b, mKiller2, board.NOMOVE); score != SCORE_KILLER2 {
			t.Errorf("Killer 2 score = %d; want %d", score, SCORE_KILLER2)
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

		mEP := board.NewMove(board.E5, board.D6, board.MoveEnPass)
		score := s.scoreMove(&b, mEP, board.NOMOVE)

		wantScore := MVVLVA[board.WP][board.BP] + SCORE_CAPTURE
		if score != wantScore {
			t.Errorf("En passant capture score = %d; want %d", score, wantScore)
		}
	})

	t.Run("score buckets are strictly ordered and history cannot outrank them", func(t *testing.T) {
		if MAX_HISTORY >= SCORE_KILLER2 {
			t.Errorf("MAX_HISTORY (%d) must be below SCORE_KILLER2 (%d)", MAX_HISTORY, SCORE_KILLER2)
		}

		lowestCapture := MVVLVA[board.WP][board.BK] 
		for _, v := range MVVLVA {
			for _, x := range v {
				if x < lowestCapture {
					lowestCapture = x
				}
			}
		}

		tiers := []struct {
			name  string
			score int
		}{
			{"TT", SCORE_TT},
			{"PV", SCORE_PV},
			{"lowest capture", SCORE_CAPTURE + lowestCapture},
			{"promotion", SCORE_PROMOTION},
			{"killer 1", SCORE_KILLER1},
			{"killer 2", SCORE_KILLER2},
			{"max history", MAX_HISTORY},
		}
		for i := 1; i < len(tiers); i++ {
			if tiers[i-1].score <= tiers[i].score {
				t.Errorf("%s (%d) must outrank %s (%d)",
					tiers[i-1].name, tiers[i-1].score, tiers[i].name, tiers[i].score)
			}
		}
	})

	t.Run("saturated history does not outrank captures or killers", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
		b.Ply = 0

		mCapture := board.NewMove(board.E4, board.D5, board.MoveCapture)
		s.stack[b.Ply].killers[0] = mKiller1
		s.History[board.WP][board.A3] = MAX_HISTORY

		var ml board.MoveList
		ml.AddMove(mHistory)
		ml.AddMove(mKiller1)
		ml.AddMove(mCapture)

		s.sortMove(&b, &ml, board.NOMOVE)

		expected := []board.Move{mCapture, mKiller1, mHistory}
		for i, want := range expected {
			if ml.Moves[i] != want {
				t.Errorf("Move at index %d = %v; want %v", i, ml.Moves[i], want)
			}
		}
	})

	t.Run("sortMove arranges moves in descending order of score", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		b.ParseFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
		b.Ply = 0

		s.stack[b.Ply].killers[0] = mKiller1
		s.History[board.WP][board.A3] = 120

		mCapture := board.NewMove(board.E4, board.D5, board.MoveCapture)

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

func TestUpdateKiller(t *testing.T) {
	m1 := board.NewMove(board.G1, board.F3, board.MoveQuiet)
	m2 := board.NewMove(board.B1, board.C3, board.MoveQuiet)

	t.Run("new killer shifts the previous one into slot 2", func(t *testing.T) {
		s := &Searcher{}
		s.updateKiller(3, m1)
		s.updateKiller(3, m2)

		if got := s.stack[3].killers; got[0] != m2 || got[1] != m1 {
			t.Errorf("killers = %v; want [%v %v]", got, m2, m1)
		}
	})

	t.Run("repeating the primary killer does not duplicate it", func(t *testing.T) {
		s := &Searcher{}
		s.updateKiller(3, m1)
		s.updateKiller(3, m2)
		s.updateKiller(3, m2) 

		if got := s.stack[3].killers; got[0] != m2 || got[1] != m1 {
			t.Errorf("killers = %v; want [%v %v] (second killer lost)", got, m2, m1)
		}
	})
}
