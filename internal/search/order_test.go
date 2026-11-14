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

		if score := s.scoreMove(&b, mTT, mTT); score != ScoreTT {
			t.Errorf("TT move score = %d; want %d", score, ScoreTT)
		}

		s.PV.ScorePV = true
		s.PV.Table[0][b.Ply] = mPV
		if score := s.scoreMove(&b, mPV, board.NOMOVE); score != ScorePV {
			t.Errorf("PV move score = %d; want %d", score, ScorePV)
		}
		if s.PV.ScorePV {
			t.Errorf("ScorePV was not reset to false after scoring PV move")
		}

		mCapture := board.NewMove(board.E4, board.D5, board.MoveCapture)
		expectedCapScore := MVVLVA[board.WP][board.BP] + ScoreCapture
		if score := s.scoreMove(&b, mCapture, board.NOMOVE); score != expectedCapScore {
			t.Errorf("Capture score = %d; want %d", score, expectedCapScore)
		}

		if score := s.scoreMove(&b, mPromoQ, board.NOMOVE); score != ScoreQueenPromo {
			t.Errorf("Queen promo score = %d; want %d", score, ScoreQueenPromo)
		}

		s.stack[b.Ply].killers[0] = mKiller1
		s.stack[b.Ply].killers[1] = mKiller2
		if score := s.scoreMove(&b, mKiller1, board.NOMOVE); score != ScoreKiller1 {
			t.Errorf("Killer 1 score = %d; want %d", score, ScoreKiller1)
		}
		if score := s.scoreMove(&b, mKiller2, board.NOMOVE); score != ScoreKiller2 {
			t.Errorf("Killer 2 score = %d; want %d", score, ScoreKiller2)
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

		wantScore := MVVLVA[board.WP][board.BP] + ScoreCapture
		if score != wantScore {
			t.Errorf("En passant capture score = %d; want %d", score, wantScore)
		}
	})
	t.Run("scoreMove ranks promotions by piece", func(t *testing.T) {
		s := &Searcher{}
		b := board.NewBoard()
		// white pawn on b7 can push to b8 or capture the rook on a8
		b.ParseFEN("r3k3/1P6/8/8/8/8/8/4K3 w - - 0 1")
		b.Ply = 0

		quietQ := board.NewMove(board.B7, board.B8, board.MovePromoQ)
		quietN := board.NewMove(board.B7, board.B8, board.MovePromoN)
		capQ := board.NewMove(board.B7, board.A8, board.MovePromoCaptureQ)

		if got := s.scoreMove(&b, quietQ, board.NOMOVE); got != ScoreQueenPromo {
			t.Errorf("quiet queen promotion = %d; want %d", got, ScoreQueenPromo)
		}

		if got := s.scoreMove(&b, quietN, board.NOMOVE); got != ScoreUnderPromo {
			t.Errorf("quiet underpromotion = %d; want %d", got, ScoreUnderPromo)
		}

		wantCapQ := MVVLVA[board.WP][board.BR] + ScoreQueenPromo
		if got := s.scoreMove(&b, capQ, board.NOMOVE); got != wantCapQ {
			t.Errorf("queen capture-promotion = %d; want %d", got, wantCapQ)
		}

		plainCapture := MVVLVA[board.WP][board.BR] + ScoreCapture
		if wantCapQ <= plainCapture {
			t.Errorf("queen capture-promotion (%d) should outrank a plain capture of the same piece (%d)", wantCapQ, plainCapture)
		}
		if wantCapQ <= ScoreQueenPromo {
			t.Errorf("queen capture-promotion (%d) should outrank a quiet queen promotion (%d)", wantCapQ, ScoreQueenPromo)
		}
	})
	t.Run("score tiers do not overlap", func(t *testing.T) {
		maxMVV, minMVV := MVVLVA[0][0], MVVLVA[0][0]
		for _, row := range MVVLVA {
			for _, v := range row {
				maxMVV = max(maxMVV, v)
				minMVV = min(minMVV, v)
			}
		}

		switch {
		case ScoreTT <= ScorePV:
			t.Error("TT tier must sit above the PV tier")
		case ScorePV <= ScoreQueenPromo+maxMVV:
			t.Error("queen capture-promotions can reach the PV tier")
		case ScoreQueenPromo <= ScoreCapture+maxMVV:
			t.Error("plain captures can reach the queen-promotion tier")
		case ScoreCapture+minMVV <= ScoreKiller1:
			t.Error("captures can fall into the killer tier")
		case ScoreKiller1 <= ScoreKiller2:
			t.Error("killer 1 must outrank killer 2")
		case ScoreKiller2 <= MAX_HISTORY:
			t.Error("history scores can reach the killer tier")
		case -MAX_HISTORY <= ScoreUnderPromo:
			t.Error("history scores can fall to the underpromotion tier")
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
