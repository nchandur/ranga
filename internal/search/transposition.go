package search

import (
	"ranga/internal/board"
	"unsafe"
)

// evaluation flags for transposition table entries
const (
	FEXACT int = iota // indicates exact evaluation value.
	FALPHA            // indicates upper bound
	FBETA             // indicates lower bound

)

const NOENTRY int = 99999 // indicate cache miss or uninitialized entry

// holds tranposition table entries
type TranspositionTableEntry struct {
	board.Move
	Score int
	Depth int
	Flag  int
	Key   uint64
}

// transposition table
// cache of previously evaluated chess positions
type TranspositionTable struct {
	Entries []TranspositionTableEntry // actual cached evaluation entries
	Length  int                       // max capacity of transposition table
}

// instantiates new transposition table
// size specifies the total memory allocated for the transposition table in MB
func NewTranspositionTable(size int) *TranspositionTable {

	size = max(size, 1)
	bytes := size * 1024 * 1024

	entrySize := int(unsafe.Sizeof(TranspositionTableEntry{}))
	targetEntries := bytes / entrySize

	entryCount := 1
	for entryCount*2 <= targetEntries {
		entryCount *= 2
	}

	return &TranspositionTable{Entries: make([]TranspositionTableEntry, entryCount), Length: entryCount}
}

// clears all entries in transposition table
func (tt *TranspositionTable) Clear() {
	for i := range tt.Length {
		tt.Entries[i] = TranspositionTableEntry{}
	}
}

// checks transposition table for previously evaluated position
func (tt *TranspositionTable) Probe(alpha, beta, ply, depth int, key uint64) int {

	idx := key % uint64(tt.Length)

	entry := (*tt).Entries[idx]

	if entry.Key == key {

		score := entry.Score

		if score < -MATESCORE {
			score += ply
		}

		if score > MATESCORE {
			score -= ply
		}

		if entry.Depth >= depth {
			if entry.Flag == FEXACT {
				return score
			}

			if (entry.Flag == FALPHA) && (score <= alpha) {
				return alpha
			}

			if (entry.Flag == FBETA) && (score >= beta) {
				return beta
			}

		}
	}

	return NOENTRY
}

// saves or updates entry in transposition table
func (tt *TranspositionTable) Store(score, depth, ply, flag int, key uint64, move board.Move) {

	entry := &(tt.Entries[key%uint64(tt.Length)])

	if entry.Key != 0 && entry.Depth > depth {
		return
	}

	if score < -MATESCORE {
		score -= ply
	}

	if score > MATESCORE {
		score += ply
	}

	entry.Key = key
	entry.Score = score
	entry.Flag = flag
	entry.Depth = depth
	entry.Move = move

}

// returns move stored in transposition table
func (tt *TranspositionTable) ProbeMove(key uint64) board.Move {

	entry := &tt.Entries[key%uint64(tt.Length)]

	if entry.Key == key {
		return entry.Move
	}

	return board.NOMOVE
}
