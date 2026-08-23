package search

import (
	"fmt"
	"ranga/internal/board"
	"strings"
)

// stores principal variation lines
type PVTable struct {
	Table    [MAX_PLY][MAX_PLY + 1]board.Move // sequence of best moves for each ply depth
	Length   [MAX_PLY + 1]int                 // number of valid principal variation moves recorded at each ply
	FollowPv bool                             // indicates whether search is currently traversing existing pv line
	ScorePV  bool                             // indicates whether pv moves should be prioritized for move ordering at current node
}

// resets the principal variation table
func (p *PVTable) Clear() {
	p.Length = [MAX_PLY + 1]int{}
	p.Table = [MAX_PLY][MAX_PLY + 1]board.Move{}
	p.FollowPv = false
	p.ScorePV = false
}

// flags pv move at given search depth to be prioritized for move ordering
func (p *PVTable) enablePVScoring(ml *board.MoveList, ply int) {

	p.FollowPv = false

	for count := range ml.Count {
		if p.Table[0][ply] == ml.Moves[count] {
			p.ScorePV = true
			p.FollowPv = true
		}
	}

}

// update line in pv at given ply
func (p *PVTable) updatePVLine(move board.Move, ply int) {
	p.Table[ply][0] = move

	for i := 0; i < p.Length[ply+1]; i++ {
		p.Table[ply][i+1] = p.Table[ply+1][i]
	}

	p.Length[ply] = p.Length[ply+1] + 1
}

func (p PVTable) String() string {

	var res strings.Builder

	for cnt := range p.Length[0] {
		fmt.Fprintf(&res, "%s ", p.Table[0][cnt])
	}

	return res.String()
}
