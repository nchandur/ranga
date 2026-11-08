package uci

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"ranga/internal/board"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	genfensPlies   = 8
	genfensEvalCap = 400  // reject openings scoring beyond +/- this many centipawns
	genfensNodes   = 5000 // node budget for scoring each candidate
)

// loads an opening book (EPD or FEN, one position per line)
func loadBook(path string) ([]string, error) {
	if path == "" || strings.EqualFold(path, "none") {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var fens []string
	skipped := 0
	scratch := board.NewBoard()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 {
			continue
		}
		if len(fields) < 4 {
			skipped++
			continue
		}

		fen := strings.Join(fields[:4], " ") + " 0 1"
		if len(fields) >= 6 {
			_, e1 := strconv.Atoi(fields[4])
			_, e2 := strconv.Atoi(fields[5])
			if e1 == nil && e2 == nil {
				fen = strings.Join(fields[:6], " ")
			}
		}

		if scratch.ParseFEN(fen) != nil {
			skipped++
			continue
		}
		fens = append(fens, fen)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	if len(fens) == 0 {
		return nil, fmt.Errorf("no valid positions in %s", path)
	}
	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "genfens: skipped %d unreadable lines in %s\n", skipped, path)
	}
	return fens, nil
}

// generates `count` random openings
func (e *Engine) generate(count int, seed uint64, book []string, emit func(*board.Board)) int {
	// scoring candidates needs a fresh, unset stop flag, as handleGo does
	e.pauseSearch()
	stop := new(atomic.Bool)
	e.stop = stop
	e.searcher.Stop = stop

	accept := func(b *board.Board) bool {
		return abs(e.scorePosition(b)) <= genfensEvalCap
	}

	n := genPositions(&e.board, book, count, genfensPlies, seed, accept, emit)

	e.searcher.NN.Reset(&e.board)
	e.searcher.TT.Clear()

	return n
}

// generates `count` positions and saves just the FENs
func (e *Engine) GenFENsToFile(path string, count int, seed uint64, bookPath string) (int, error) {
	book, err := loadBook(bookPath)
	if err != nil {
		return 0, fmt.Errorf("book: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	w := bufio.NewWriter(f)

	n := e.generate(count, seed, book, func(b *board.Board) {
		fmt.Fprintln(w, b.FEN())
	})

	if err := w.Flush(); err != nil {
		f.Close()
		return n, err
	}
	if err := f.Close(); err != nil {
		return n, err
	}
	if n < count {
		return n, fmt.Errorf("produced %d of %d positions", n, count)
	}
	return n, nil
}

// plays `plies` random legal moves from a starting position, `count` times
func genPositions(start *board.Board, book []string, count, plies int, seed uint64,
	accept func(*board.Board) bool, emit func(*board.Board)) int {

	rng := rand.New(rand.NewPCG(seed, seed))

	produced := 0
	for attempts := 0; produced < count && attempts < 1000*count; attempts++ {
		b := start.Preserve()
		if len(book) > 0 {
			b = board.NewBoard()
			if b.ParseFEN(book[rng.IntN(len(book))]) != nil {
				continue
			}
		}

		ok := true
		for range plies {
			if !playRandomMove(&b, rng) {
				ok = false
				break
			}
		}
		if !ok || !hasLegalMove(&b) {
			continue
		}
		if accept != nil && !accept(&b) {
			continue
		}

		emit(&b)
		produced++
	}
	return produced
}

// plays one random legal move
func playRandomMove(b *board.Board, rng *rand.Rand) bool {
	list := board.NewMoveList()
	list.GenerateMoves(b)

	for _, i := range rng.Perm(list.Count) {
		if b.MakeMove(list.Moves[i], false) {
			b.Ply++
			return true
		}
	}
	return false
}

// reports whether side to move has at least one legal move
func hasLegalMove(b *board.Board) bool {
	list := board.NewMoveList()
	list.GenerateMoves(b)

	for i := range list.Count {
		c := b.Preserve()
		if c.MakeMove(list.Moves[i], false) {
			return true
		}
	}
	return false
}
