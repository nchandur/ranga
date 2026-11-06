package search

import "ranga/internal/board"

type Node struct {
	key        uint64
	staticEval int           // static eval at this ply
	killers    [2]board.Move // holds killer moves
	excluded   board.Move    // move to skip at each ply
}

// headroom so stack[ply+N] lookups never go out of bounds
type NodeStack [MAX_PLY + 8]Node

func (st *NodeStack) reset() {
	for i := range st {
		st[i] = Node{
			staticEval: -INFINITY,
			killers:    [2]board.Move{board.NOMOVE, board.NOMOVE},
			excluded:   board.NOMOVE,
		}
	}
}

func (s *Searcher) pushKiller(ply int, m board.Move) {
	n := &s.stack[ply]
	if n.killers[0] != m {
		n.killers[1] = n.killers[0]
		n.killers[0] = m
	}
}
