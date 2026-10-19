package search

import (
	"math/bits"
	"ranga/internal/board"
	"sync/atomic"
	"unsafe"
)

// evaluation flags for transposition table entries
const (
	FEXACT int = iota // indicates exact evaluation value
	FALPHA            // indicates upper bound
	FBETA             // indicates lower bound

)

const NOENTRY int = 99999 // indicate cache miss or uninitialized entry

// A TT entry is packed into one uint64:
//
//	bits  0-15  move                      (16 bits, board.Move as is)
//	bits 16-37  score + scoreBias         (22 bits, range about +-2.09M)
//	bits 38-45  depth                     (8 bits, 0..255)
//	bits 46-47  flag + 1                  (2 bits, never 0)
//	bits 48-63  key check                 (top 16 bits of the zobrist key)

const (
	moveBits  = 16
	scoreBits = 22
	depthBits = 8
	flagBits  = 2

	moveShift  = 0
	scoreShift = moveShift + moveBits
	depthShift = scoreShift + scoreBits
	flagShift  = depthShift + depthBits
	checkShift = flagShift + flagBits // 48

	moveMask  = 1<<moveBits - 1
	scoreMask = 1<<scoreBits - 1
	depthMask = 1<<depthBits - 1
	flagMask  = 1<<flagBits - 1

	scoreBias = 1 << (scoreBits - 1)
)

type TranspositionTableEntry uint64

func packEntry(key uint64, score, depth, flag int, move board.Move) TranspositionTableEntry {
	score = max(-scoreBias+1, min(scoreBias-1, score))
	depth = max(0, min(depthMask, depth))

	return TranspositionTableEntry(
		uint64(move)<<moveShift |
			uint64(score+scoreBias)<<scoreShift |
			uint64(depth)<<depthShift |
			uint64(flag+1)<<flagShift |
			(key>>checkShift)<<checkShift)
}

func (e TranspositionTableEntry) Move() board.Move {
	return board.Move((uint64(e) >> moveShift) & moveMask)
}
func (e TranspositionTableEntry) Score() int {
	return int((uint64(e)>>scoreShift)&scoreMask) - scoreBias
}
func (e TranspositionTableEntry) Depth() int {
	return int((uint64(e) >> depthShift) & depthMask)
}
func (e TranspositionTableEntry) Flag() int {
	return int((uint64(e)>>flagShift)&flagMask) - 1
}

// reports whether entry is occupied and belongs to key
func (e TranspositionTableEntry) matches(key uint64) bool {
	return e != 0 && uint64(e)>>checkShift == key>>checkShift
}

// cache of previously evaluated chess positions
type TranspositionTable struct {
	Entries []atomic.Uint64 // packed entries
	Length  int             // max capacity of transposition table (power of two)
}

// returns table slot for key
func (tt *TranspositionTable) slot(key uint64) *atomic.Uint64 {
	return &tt.Entries[key&uint64(tt.Length-1)]
}

// instantiates new transposition table
func NewTranspositionTable(size int) *TranspositionTable {
	tt := &TranspositionTable{}
	tt.Resize(size)
	return tt
}

// clears all entries in transposition table.
func (tt *TranspositionTable) Clear() {
	clear(tt.Entries)
}

// checks transposition table for previously evaluated position
func (tt *TranspositionTable) Probe(alpha, beta, ply, depth int, key uint64) int {

	entry := TranspositionTableEntry(tt.slot(key).Load())

	if entry.matches(key) {

		score := entry.Score()

		if score < -MATESCORE {
			score += ply
		}

		if score > MATESCORE {
			score -= ply
		}

		if entry.Depth() >= depth {
			flag := entry.Flag()

			if flag == FEXACT {
				return score
			}

			if (flag == FALPHA) && (score <= alpha) {
				return alpha
			}

			if (flag == FBETA) && (score >= beta) {
				return beta
			}

		}
	}

	return NOENTRY
}

// saves or updates entry in transposition table
func (tt *TranspositionTable) Store(score, depth, ply, flag int, key uint64, move board.Move) {

	slot := tt.slot(key)

	old := TranspositionTableEntry(slot.Load())
	if old != 0 && old.Depth() > depth {
		return
	}

	if score < -MATESCORE {
		score -= ply
	}

	if score > MATESCORE {
		score += ply
	}

	slot.Store(uint64(packEntry(key, score, depth, flag, move)))
}

// returns move stored in transposition table
func (tt *TranspositionTable) ProbeMove(key uint64) board.Move {

	entry := TranspositionTableEntry(tt.slot(key).Load())

	if entry.matches(key) {
		return entry.Move()
	}

	return board.NOMOVE
}

// allocates a new TT based on the requested size in MB
func (tt *TranspositionTable) Resize(size int) {
	size = max(size, 1)
	sizeInBytes := uint64(size) * 1024 * 1024

	entrySize := uint64(unsafe.Sizeof(atomic.Uint64{}))
	maxEntries := sizeInBytes / entrySize

	powerOfTwoIndex := bits.Len64(maxEntries) - 1
	numEntries := int(uint64(1) << powerOfTwoIndex)

	tt.Entries = make([]atomic.Uint64, numEntries)
	tt.Length = numEntries
}

// returns full entry for key with score adjusted for mate distance at ply
func (tt *TranspositionTable) ProbeEntry(key uint64, ply int) (move board.Move, score, depth, flag int, ok bool) {
	entry := TranspositionTableEntry(tt.slot(key).Load())

	if !entry.matches(key) {
		return board.NOMOVE, 0, 0, 0, false
	}

	score = entry.Score()

	if score < -MATESCORE {
		score += ply
	}

	if score > MATESCORE {
		score -= ply
	}

	return entry.Move(), score, entry.Depth(), entry.Flag(), true
}
