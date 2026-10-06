package uci

import (
	"context"
	"fmt"
	"os"
	"ranga/internal/board"
	"ranga/internal/search"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// handles uci command
// introduces engine
func (e *Engine) handleUCI() {
	e.writeLine(fmt.Sprintf("id name ranga %s", e.version))
	e.writeLine("id author nchandur")
	e.options.Print()
	e.writeLine("uciok")
}

// handles setoption command
func (e *Engine) handleSetOption(command []string) {
	payload := strings.Join(command, " ")
	e.options.Set(payload)
}

// handles quit command
// quits main loop and exits
func (e *Engine) handleQuit() {
	e.pauseSearch()
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
	e.pauseSearch()
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

	// reset nn
	e.searcher.NN.Reset(&e.board)
}

// handles go command
func (e *Engine) handleGo(args []string) {
	e.pauseSearch()

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
		case "nodes":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					opts.nodes = val
				}
				i++
			}
		case "perft":
			if i+1 < len(args) {
				if d, err := strconv.Atoi(args[i+1]); err == nil {
					opts.depth = d
					opts.perft = true
				}
				i++
			}
		}
	}
	var timeAllocation TimeAllocation
	if !opts.infinite && !opts.perft {
		timeAllocation = e.calculateTimeLimit(opts)
	}

	if opts.perft {
		e.runPerft(context.Background(), opts.depth)
		return
	}

	stop := new(atomic.Bool)
	e.stop = stop
	e.searcher.Stop = stop

	var timer *time.Timer
	if timeAllocation.Hard > 0 {
		timer = time.AfterFunc(timeAllocation.Hard, func() { stop.Store(true) })
	}

	e.searchWg.Go(func() {
		if timer != nil {
			defer timer.Stop()
		}
		defer stop.Store(true)
		e.runSearch(stop, opts, timeAllocation)
	})
}

// helper function to run perft on position
func (e *Engine) runPerft(ctx context.Context, depth int) {
	board.PerftDivide(ctx, &e.board, depth)
}

// helper function to run search and evaluation
func (e *Engine) runSearch(stop *atomic.Bool, opts goOptions, timeAllocation TimeAllocation) {
	maxDepth := search.MAX_DEPTH
	if opts.depth > 0 && !opts.infinite {
		maxDepth = opts.depth
	}

	e.searcher.NodeLimit = opts.nodes
	e.searcher.Nodes = 0
	e.searcher.NodesPub.Store(0)

	// launch helpers
	var helpersWg sync.WaitGroup
	for i, h := range e.helpers {
		h.Stop = stop
		h.NodeLimit = 0
		h.Nodes = 0
		h.NodesPub.Store(0)

		hb := e.board
		helpersWg.Go(func() { runHelper(h, hb, i+1, maxDepth, stop) })
	}

	bestMove, prevBestMove := board.NOMOVE, board.NOMOVE
	stableIterations := 0
	previousScore := 0
	firstScore := true

	searchStart := time.Now()

	for d := 1; d <= maxDepth; d++ {
		if timeAllocation.Soft > 0 && time.Since(searchStart) > timeAllocation.Soft {
			break
		}

		move, score := e.searcher.Search(&e.board, d)

		if stop.Load() {
			break
		}

		if move != board.NOMOVE {
			bestMove = move
			e.writeLine(fmt.Sprintf("info depth %d score cp %d nodes %d pv %s", d, score, e.totalNodes(), e.searcher.PV))
		}

		// bank unused time once search has settled
		if bestMove == prevBestMove {
			stableIterations++
		} else {
			stableIterations = 0
		}
		prevBestMove = bestMove

		if !firstScore && timeAllocation.Soft > 0 {
			diff := score - previousScore

			// score dropped sharply
			if diff < -100 {
				timeAllocation.Soft = timeAllocation.Hard
			} else if abs(diff) > 50 {
				// general volatility, capped at hard limit
				extended := min(timeAllocation.Soft*3/2, timeAllocation.Hard)
				timeAllocation.Soft = extended
			}
		}
		previousScore = score
		firstScore = false

		if timeAllocation.Soft > 0 && d >= 6 && stableIterations >= 4 && time.Since(searchStart) > timeAllocation.Soft/3 {
			break
		}
	}

	// thread 0 is done (depth, time, nodes or stop): halt helpers and wait for them
	// before touching e.board or printing bestmove
	stop.Store(true)
	helpersWg.Wait()

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
	e.stop.Store(true)
	e.searchWg.Wait()
}

func (e *Engine) totalNodes() int {
	n := e.searcher.Nodes
	for _, h := range e.helpers {
		n += int(h.NodesPub.Load())
	}
	return n
}

func runHelper(h *search.Searcher, b board.Board, id, maxDepth int, stop *atomic.Bool) {
	for d := 1; d <= maxDepth && !stop.Load(); d++ {
		if skipDepth(id, d) {
			continue
		}
		h.Search(&b, d)
	}
}

var skipSize = [20]int{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4}
var skipPhase = [20]int{0, 1, 0, 1, 2, 3, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7}

func skipDepth(id, d int) bool {
	i := (id - 1) % 20
	return ((d+skipPhase[i])/skipSize[i])%2 != 0
}
