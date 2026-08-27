package uci

import (
	"context"
	"fmt"
	"os"
	"ranga/internal/board"
	"ranga/internal/search"
	"strconv"
	"strings"
)

// handles uci command
// introduces engine
func (e *Engine) handleUCI() {
	e.writeLine(fmt.Sprintf("id name ranga %s", e.version))
	e.writeLine("id author nchandur")
	e.writeLine("uciok")
}

// handles quit command
// quits main loop and exits
func (e *Engine) handleQuit() {
	if e.searchCancel != nil {
		e.searchCancel()
	}

	e.searchWg.Wait()
}

// handles show command
// displays pieces on an ASCII chessboard
func (e *Engine) handleShow() {
	e.pauseSearch()
	e.board.Print()
}

// handles isready command
// confirms if engine can send and receive info
func (e *Engine) handleIsReady() {
	e.writeLine("readyok")
}

// handle ucinewgame command
// sets board to new game
func (e *Engine) handleNewGame() {
	e.handlePosition([]string{"position", "startpos"})
}

// handles clear command
// removes all pieces from board
func (e *Engine) handleClear() {
	e.pauseSearch()
	e.board.Clear()
}

// handles stop command
// stops searching tree and returns current best move
func (e *Engine) handleStop() {
	if e.searchCancel != nil {
		e.searchCancel()
	}
	e.searchWg.Wait()
	e.searchCancel = nil
}

// handles eval command
// provides an evaluation of current position on board
func (e *Engine) handleEvaluate() {
	e.pauseSearch()
	e.writeLine(fmt.Sprintf("score %.2f", float64(e.searcher.Evaluate(&e.board))/float64(100)))
}

// handles position command
// sets pieces on board after parsing valid FEN string
func (e *Engine) handlePosition(args []string) {

	e.pauseSearch()

	line := strings.Join(args, " ")

	var fenStr string
	var movesStr string

	// sets up starting position
	if strings.HasPrefix(line, "startpos") {
		fenStr = board.START
		if _, after, found := strings.Cut(line, "moves"); found {
			movesStr = after
		}
		// sets up position based on valid FEN string
	} else if after, ok := strings.CutPrefix(line, "fen"); ok {
		remaining := strings.TrimSpace(after)

		before, after, found := strings.Cut(remaining, "moves")
		if found {
			fenStr = strings.TrimSpace(before)
			movesStr = after
		} else {
			fenStr = remaining
		}
		// sets up starting position if invalid fen string is passed
	} else {
		fenStr = board.START
	}

	if err := e.board.ParseFEN(fenStr); err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing FEN '%s': %v\n", fenStr, err)
		return
	}

	e.board.Repetition.Idx = 0
	e.board.Repetition.Table[e.board.Repetition.Idx] = e.board.Key

	movesStr = strings.TrimSpace(movesStr)
	if movesStr != "" {
		for moveStr := range strings.FieldsSeq(movesStr) {
			move := e.board.ParseMove(moveStr)

			if move == 0 {
				break
			}

			ok := e.board.MakeMove(move, false)
			if !ok {
				fmt.Fprintf(os.Stderr, "illegal move encountered in setup: %s\n", moveStr)
				return
			}

			e.board.Repetition.Idx++
			e.board.Repetition.Table[e.board.Repetition.Idx] = e.board.Key

		}
	}

	e.board.Ply = 0
}

// handles go command
func (e *Engine) handleGo(args []string) {
	if e.searchCancel != nil {
		e.handleStop()
	}

	opts := goOptions{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "depth":
			if i+1 < len(args) {
				if d, err := strconv.Atoi(args[i+1]); err == nil {
					opts.depth = d
				}
				i++
			}
		case "infinite":
			opts.infinite = true
		case "wtime":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.wtime = val
				}
				i++
			}
		case "btime":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.btime = val
				}
				i++
			}
		case "winc":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.winc = val
				}
				i++
			}
		case "binc":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.binc = val
				}
				i++
			}
		case "movestogo":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.movesToGo = val
				}
				i++
			}
		case "movetime":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.moveTime = val
				}
				i++
			}
		case "metrics":
			if i+1 < len(args) {
				if d, err := strconv.Atoi(args[i+1]); err == nil {
					opts.depth = d
				}
				i++
			}
		}
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if !opts.infinite && !opts.metrics {
		if timeAllocation := e.calculateTimeLimit(opts); timeAllocation > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), timeAllocation)
		}
	}

	if ctx == nil {
		ctx, cancel = context.WithCancel(context.Background())
	}
	e.searchCancel = cancel

	e.searchWg.Go(func() {
		defer cancel()
		e.runSearch(ctx, opts)
	})
}

// helper function to run search and evaluation
func (e *Engine) runSearch(ctx context.Context, opts goOptions) {
	maxDepth := search.MAX_DEPTH
	if opts.depth > 0 && !opts.infinite {
		maxDepth = opts.depth
	}

	bestMove := board.NOMOVE

	for d := 1; d <= maxDepth; d++ {
		move, score, nodes := e.searcher.Search(ctx, &e.board, d)

		if ctx.Err() != nil {
			break
		}

		if move != board.NOMOVE {
			bestMove = move
			e.writeLine(fmt.Sprintf("info depth %d score cp %d nodes %d pv %s", d, score, nodes, e.searcher.PV))
		}

	}

	if bestMove == board.NOMOVE {
		ml := board.NewMoveList()
		ml.GenerateMoves(&e.board)
		for _, m := range ml.Moves[:ml.Count] {
			state := e.board.Preserve()
			if e.board.MakeMove(m, false) {
				e.board.Restore(&state)
				bestMove = m
				break
			}
			e.board.Restore(&state)
		}
	}

	e.writeLine("bestmove " + bestMove.String())
}

// helper function to pause search
func (e *Engine) pauseSearch() {
	if e.searchCancel != nil {
		e.searchCancel()
		e.searchWg.Wait()
		e.searchCancel = nil
	}
}
