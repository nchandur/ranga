package search

import "ranga/internal/board"

type Node struct {
	key        uint64
	staticEval int           // static eval at this ply
	killers    [2]board.Move // holds killer moves
	excluded   board.Move    // move to skip at each ply
}

type NodeStack [MAX_PLY + 8]Node

func (ns *NodeStack) reset() {
	for i := range ns {
		ns[i] = Node{
			staticEval: -INFINITY,
			killers:    [2]board.Move{board.NOMOVE, board.NOMOVE},
			excluded:   board.NOMOVE,
		}
	}
}
