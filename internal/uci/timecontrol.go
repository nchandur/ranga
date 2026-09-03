package uci

import (
	"ranga/internal/board"
	"time"
)

const MoveOverhead int = 50

// manages options for go command
type goOptions struct {
	depth     int
	infinite  bool
	metrics   bool
	wtime     int
	btime     int
	winc      int
	binc      int
	movesToGo int
	moveTime  int
	nodes     int
}

type TimeAllocation struct {
	Soft time.Duration // checked between iterations
	Hard time.Duration // absolute cutoff
}

// helper function to manage time during games
func (e *Engine) calculateTimeLimit(opts goOptions) TimeAllocation {
	if opts.moveTime > 0 {
		allocatedMs := opts.moveTime - MoveOverhead
		if allocatedMs < 10 {
			allocatedMs = 10
		}
		d := time.Duration(allocatedMs) * time.Millisecond
		return TimeAllocation{Soft: d, Hard: d}
	}

	var timeLeft, increment int
	if e.board.Side == board.White {
		timeLeft = opts.wtime
		increment = opts.winc
	} else {
		timeLeft = opts.btime
		increment = opts.binc
	}

	if timeLeft <= 0 {
		return TimeAllocation{}
	}

	movesToGo := 40
	if opts.movesToGo > 0 {
		movesToGo = opts.movesToGo
	}

	base := (timeLeft / movesToGo) + (increment * 3 / 4) - MoveOverhead

	softMs := base
	hardMs := base * 3 // allow overrun for a single hard iteration, bounded below

	maxAllowed := timeLeft * 4 / 5
	if hardMs > maxAllowed {
		hardMs = maxAllowed
	}
	if softMs > hardMs {
		softMs = hardMs
	}

	if softMs < 10 {
		softMs = 10
	}
	if hardMs < 10 {
		hardMs = 10
	}

	return TimeAllocation{
		Soft: time.Duration(softMs) * time.Millisecond,
		Hard: time.Duration(hardMs) * time.Millisecond,
	}
}
