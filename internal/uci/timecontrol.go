package uci

import (
	"ranga/internal/board"
	"time"
)

// computes how long to search for, given UCI go options and side to move
func calculateSearchTime(opts goOptions, side board.Color) time.Duration {
	// explicit movetime overrides everything
	if opts.moveTime > 0 {
		return time.Duration(opts.moveTime) * time.Millisecond
	}

	var timeLeft, inc int
	if side == board.White {
		timeLeft, inc = opts.wtime, opts.winc
	} else {
		timeLeft, inc = opts.btime, opts.binc
	}

	// no time info at all (e.g. infinite analysis) — caller handles separately
	if timeLeft == 0 {
		return 0
	}

	movesToGo := 30 // crude fixed default
	if opts.movesToGo > 0 {
		movesToGo = opts.movesToGo
	}

	const safetyMarginMs = 50

	allocated := timeLeft/movesToGo + inc
	allocated -= safetyMarginMs

	if allocated < 10 {
		allocated = 10
	}

	return time.Duration(allocated) * time.Millisecond
}
