package search

import (
	"sync"
	"sync/atomic"

	"ranga/internal/board"
	"ranga/internal/evaluate/nnue"
)

// state of the helper threads
type smpState struct {
	helpers []*Searcher
	abort   atomic.Bool // set to stop helpers at end of a root search
	wg      sync.WaitGroup
}

// sets number of search threads
func (s *Searcher) SetThreads(n int) {
	if s.NN == nil {
		return
	}

	want := min(max(n, 1), MAX_THREADS) - 1
	helpers := s.smp.helpers

	if want < len(helpers) {
		clear(helpers[want:])
		helpers = helpers[:want]
	}

	for len(helpers) < want {
		nn := &nnue.NNUE{Network: s.NN.Network}
		helpers = append(helpers, NewSearcher(nn, s.TT, &s.smp.abort))
	}

	s.smp.helpers = helpers
}

// number of search threads, including this one
func (s *Searcher) Threads() int {
	return len(s.smp.helpers) + 1
}

// clears state of the helper threads along with main thread's
func (s *Searcher) resetHelpers() {
	for _, h := range s.smp.helpers {
		h.Reset()
		h.Nodes = 0
	}
}

// nodes searched by this searcher and its helpers
func (s *Searcher) TotalNodes() int {
	total := s.Nodes
	for _, h := range s.smp.helpers {
		total += h.Nodes
	}
	return total
}

// depth that helper number i (0-based) searches
func helperDepth(i, depth int) int {
	if i%2 == 0 {
		return depth + 1
	}
	return depth
}

// starts the helper threads on a copy of the position
func (s *Searcher) startHelpers(b *board.Board, depth, alpha, beta int) {
	if len(s.smp.helpers) == 0 || depth < 2 {
		return
	}

	s.smp.abort.Store(false)

	for i, h := range s.smp.helpers {
		h.TT = s.TT
		h.History = s.History
		h.SetGame(s.game)

		root := *b
		hd := helperDepth(i, depth)

		s.smp.wg.Go(func() {
			h.searchHelper(root, hd, alpha, beta)
		})

	}
}

// aborts helpers and waits until every one of them has stopped
func (s *Searcher) stopHelpers() {
	s.smp.abort.Store(true)
	s.smp.wg.Wait()
}

// root search of a helper thread
func (s *Searcher) searchHelper(root board.Board, depth, alpha, beta int) {
	b := root
	s.NN.Reset(&b)
	s.path[0] = b.Key

	s.searchRoot(&b, depth, alpha, beta)
}
