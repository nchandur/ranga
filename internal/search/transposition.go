package search

import (
	"math/bits"
	"ranga/internal/board"
	"sync/atomic"
	"unsafe"
)

// evaluation flags for transposition table entries
const (
	FEXACT int = iota // indicates exact evaluation value.
	FALPHA            // indicates upper bound
	FBETA             // indicates lower bound

)

// bits 0-31 move | 32-39 depth | 40-41 flag | 42-63 score (signed, 22 bits)
const (
	depthShift = 32
	flagShift  = 40
	scoreShift = 42
	scoreMask  = 1<<22 - 1
)

const NOENTRY int = 99999 // indicate cache miss or uninitialized entry

// holds tranposition table entries
type TranspositionTableEntry struct {
	keyXor atomic.Uint64
	data   atomic.Uint64
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
	d, ok := tt.load(key)
	if !ok {
		return NOENTRY
	}

	score := d.score()

	if score < -MATESCORE {
		score += ply
	}

	if score > MATESCORE {
		score -= ply
	}

	if d.depth() >= depth {
		switch flag := d.flag(); {
		case flag == FEXACT:
			return score
		case flag == FALPHA && score <= alpha:
			return alpha
		case flag == FBETA && score >= beta:
			return beta
		}
	}

	return NOENTRY
}

// saves or updates entry in transposition table
func (tt *TranspositionTable) Store(score, depth, ply, flag int, key uint64, move board.Move) {
	e := &tt.Entries[key%uint64(tt.Length)]

	oldD := e.data.Load()
	oldKey := e.keyXor.Load() ^ oldD
	if oldKey != 0 && ttData(oldD).depth() > depth {
		return
	}

	if score < -MATESCORE {
		score -= ply
	}
	if score > MATESCORE {
		score += ply
	}

	d := uint64(packData(move, score, depth, flag))
	e.data.Store(d)
	e.keyXor.Store(key ^ d)
}

// returns move stored in transposition table
// returns move stored in transposition table
func (tt *TranspositionTable) ProbeMove(key uint64) board.Move {
	d, ok := tt.load(key)
	if !ok {
		return board.NOMOVE
	}
	return d.move()
}

// allocates a new TT based on the requested size in MB
func (tt *TranspositionTable) Resize(size int) {
	size = max(size, 1)
	sizeInBytes := uint64(size) * 1024 * 1024

	entrySize := uint64(unsafe.Sizeof(TranspositionTableEntry{}))
	maxEntries := sizeInBytes / entrySize

	powerOfTwoIndex := bits.Len64(maxEntries) - 1
	numEntries := int(uint64(1) << powerOfTwoIndex)

	tt.Entries = make([]TranspositionTableEntry, numEntries)
	tt.Length = numEntries
}

// returns the entry's data if the slot holds this key
func (tt *TranspositionTable) load(key uint64) (ttData, bool) {
	e := &tt.Entries[key%uint64(tt.Length)]
	d := e.data.Load()
	k := e.keyXor.Load()
	return ttData(d), k^d == key
}

type ttData uint64

func packData(move board.Move, score, depth, flag int) ttData {
	return ttData(uint64(uint32(move)) |
		uint64(uint8(depth))<<depthShift |
		uint64(flag&3)<<flagShift |
		(uint64(score)&scoreMask)<<scoreShift)
}

func (d ttData) move() board.Move { return board.Move(uint32(d)) }
func (d ttData) depth() int       { return int(uint8(d >> depthShift)) }
func (d ttData) flag() int        { return int(d>>flagShift) & 3 }
func (d ttData) score() int       { return int(int64(d) >> scoreShift) }
