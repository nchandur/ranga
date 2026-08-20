package uci

import (
	"ranga/internal/board"
	"time"
)

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
}

// helper function to manage time during games
func (e *Engine) calculateTimeLimit(opts goOptions) time.Duration {
	if opts.moveTime > 0 {
		return time.Duration(opts.moveTime) * time.Millisecond
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
		return 0
	}

	movesToGo := 40
	if opts.movesToGo > 0 {
		movesToGo = opts.movesToGo
	}

	allocatedMs := (timeLeft / movesToGo) + (increment * 3 / 4)

	if allocatedMs > (timeLeft * 4 / 5) {
		allocatedMs = timeLeft * 4 / 5
	}

	if allocatedMs < 10 {
		allocatedMs = 10
	}

	return time.Duration(allocatedMs) * time.Millisecond
}
