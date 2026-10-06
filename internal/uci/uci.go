package uci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"ranga/internal/board"
	"ranga/internal/evaluate/nnue"
	"ranga/internal/search"
	"strconv"
	"strings"
	"sync"
)

// manages functionality of each UCI command registered
type Handler func(args []string)

// holds all necessary fields for uci, board, search and evaluation
type Engine struct {
	in       *bufio.Scanner
	out      io.Writer
	writeMux sync.Mutex

	searchCancel context.CancelFunc
	searchWg     sync.WaitGroup

	board    board.Board
	searcher *search.Searcher
	commands map[string]Handler
	options  *OptionRegistry
	version  string
}

// instantiates new engine
func NewEngine(in io.Reader, out io.Writer, version string) *Engine {

	e := &Engine{
		in:       bufio.NewScanner(in),
		out:      out,
		commands: make(map[string]Handler),
		options:  NewOptionRegistry(),
		version:  version,
	}

	e.board = board.NewBoard()

	net, err := nnue.LoadEmbedded()

	if err != nil {
		log.Fatalf("loading network: %v", err)
	}
	nn := &nnue.NNUE{Network: *net}

	nn.Reset(&e.board)

	tt := search.NewTranspositionTable(search.DEFAULT_TT_SIZE)
	e.searcher = search.NewSearcher(nn, tt)

	e.options.Register(Option{
		Name:       "Threads",
		Type:       "spin",
		DefaultVal: "1",
		Min:        1,
		Max:        1,
		Apply:      nil,
	})

	e.options.Register(Option{
		Name:       "Hash",
		Type:       "spin",
		DefaultVal: fmt.Sprintf("%d", search.DEFAULT_TT_SIZE),
		Min:        1,
		Max:        2048,
		Apply: func(value string) error {
			if v, err := strconv.Atoi(value); err == nil {
				v := min(v, 2048)
				e.pauseSearch()
				e.searcher.TT.Resize(v)
			}
			return nil
		},
	})

	e.options.Register(Option{
		Name: "Clear Hash",
		Type: "button",
		Apply: func(value string) error {
			if value == "true" {
				e.pauseSearch()
				e.searcher.TT.Clear()
			}
			return nil
		},
	})

	e.registerCommands()

	return e
}

// registers commands for uci protocol
func (e *Engine) registerCommands() {
	e.commands["uci"] = func([]string) {
		e.handleUCI()
	}

	e.commands["isready"] = func([]string) {
		e.handleIsReady()
	}

	e.commands["setoption"] = e.handleSetOption

	e.commands["show"] = func([]string) {
		e.handleShow()
	}

	e.commands["clear"] = func([]string) {
		e.handleClear()
	}

	e.commands["ucinewgame"] = func([]string) {
		e.handleNewGame()
	}

	e.commands["position"] = e.handlePosition
	e.commands["go"] = e.handleGo

	e.commands["eval"] = func([]string) {
		e.handleEvaluate()
	}

	e.commands["stop"] = func([]string) {
		e.handleStop()
	}

	e.commands["quit"] = func([]string) {
		e.handleQuit()
	}
}

// main scanning loop
func (e *Engine) Run() {
	for e.in.Scan() {
		line := strings.TrimSpace(e.in.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		cmd := fields[0]
		args := fields[1:]

		if cmd == "quit" {
			e.handleQuit()
			return
		}

		if handler, exists := e.commands[cmd]; exists {
			handler(args)
		}
	}
}

// helper to log output
func (e *Engine) writeLine(line string) {
	e.writeMux.Lock()
	defer e.writeMux.Unlock()
	fmt.Fprintln(e.out, line)
}
