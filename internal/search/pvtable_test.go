package search

import (
	"ranga/internal/board"
	"strings"
	"testing"
)

func TestPVTable(t *testing.T) {
	m1 := board.NewMove(board.E2, board.E4, board.WP, board.Empty, false, false, false, false)
	m2 := board.NewMove(board.E7, board.E5, board.BP, board.Empty, false, false, false, false)
	m3 := board.NewMove(board.G1, board.F3, board.WN, board.Empty, false, false, false, false)

	t.Run("Clear resets all table state and flags", func(t *testing.T) {
		var pv PVTable
		pv.Table[0][0] = m1
		pv.Length[0] = 1
		pv.FollowPv = true
		pv.ScorePV = true

		pv.Clear()

		if pv.FollowPv || pv.ScorePV {
			t.Errorf("flags not cleared: FollowPv=%v, ScorePV=%v", pv.FollowPv, pv.ScorePV)
		}
		for ply := 0; ply <= MAX_PLY; ply++ {
			if pv.Length[ply] != 0 {
				t.Fatalf("Length[%d] = %d; want 0", ply, pv.Length[ply])
			}
			for i := 0; i <= MAX_PLY; i++ {
				if pv.Table[ply][i] != 0 {
					t.Fatalf("Table[%d][%d] not zeroed", ply, i)
				}
			}
		}
	})
	t.Run("updatePVLine copies child line and sets correct length", func(t *testing.T) {
		var pv PVTable
		pv.updatePVLine(m3, 2)
		if pv.Length[2] != 1 {
			t.Errorf("Length[2] = %d; want 1", pv.Length[2])
		}
		if pv.Table[2][0] != m3 {
			t.Errorf("Table[2][0] = %v; want %v", pv.Table[2][0], m3)
		}

		pv.updatePVLine(m2, 1)
		if pv.Length[1] != 2 {
			t.Errorf("Length[1] = %d; want 2", pv.Length[1])
		}
		if pv.Table[1][0] != m2 || pv.Table[1][1] != m3 {
			t.Errorf("Table[1] line mismatch: got [%v, %v]", pv.Table[1][0], pv.Table[1][1])
		}

		pv.updatePVLine(m1, 0)
		if pv.Length[0] != 3 {
			t.Errorf("Length[0] = %d; want 3", pv.Length[0])
		}
		expectedRoot := [3]board.Move{m1, m2, m3}
		for i, want := range expectedRoot {
			if pv.Table[0][i] != want {
				t.Errorf("Table[0][%d] = %v; want %v", i, pv.Table[0][i], want)
			}
		}
	})
	t.Run("enablePVScoring sets flags when pv move matches MoveList", func(t *testing.T) {
		var pv PVTable
		pv.Table[0][0] = m1

		var ml board.MoveList
		ml.AddMove(m2)
		ml.AddMove(m1)

		pv.enablePVScoring(&ml, 0)

		if !pv.ScorePV || !pv.FollowPv {
			t.Errorf("expected ScorePV and FollowPv to be true, got ScorePV=%v, FollowPv=%v",
				pv.ScorePV, pv.FollowPv)
		}
	})
	t.Run("enablePVScoring remains false when pv move is absent", func(t *testing.T) {
		var pv PVTable
		pv.Table[0][0] = m1

		var ml board.MoveList
		ml.AddMove(m2)
		ml.AddMove(m3)

		pv.FollowPv = true
		pv.ScorePV = true

		pv.enablePVScoring(&ml, 0)

		if pv.ScorePV || pv.FollowPv {
			t.Errorf("expected ScorePV and FollowPv to be false, got ScorePV=%v, FollowPv=%v",
				pv.ScorePV, pv.FollowPv)
		}
	})
	t.Run("String formats root variation properly", func(t *testing.T) {
		var pv PVTable
		pv.Table[0][0] = m1
		pv.Table[0][1] = m2
		pv.Length[0] = 2

		got := strings.TrimSpace(pv.String())
		want := strings.TrimSpace(m1.String() + " " + m2.String())

		if got != want {
			t.Errorf("String() mismatch:\ngot:  %q\nwant: %q", got, want)
		}
	})
}
