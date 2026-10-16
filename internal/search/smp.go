package search

import (
	"sync"

	"ranga/internal/board"
	"ranga/internal/evaluate/nnue"
)

// reports whether helper number id skips iteration depth
func skipDepth(id, depth int) bool {
	i := (id - 1) % len(skipSize)
	return ((depth+skipPhase[i])/skipSize[i])%2 == 1
}

// publishes this searcher's node count for the main thread and reports whether the search must stop
func (s *Searcher) poll() bool {
	s.published.Store(int64(s.Nodes))

	if s.NodeLimit > 0 {
		lead := s
		if s.lead != nil {
			lead = s.lead
		}

		if len(lead.helpers) > 0 && lead.publishedNodes() >= s.NodeLimit {
			s.Stop.Store(true)
		}
	}

	return s.Stop.Load()
}

// sum of published node counts of this searcher and its helpers
func (s *Searcher) publishedNodes() int {
	total := int(s.published.Load())
	for _, h := range s.helpers {
		total += int(h.published.Load())
	}
	return total
}

// nodes searched by this searcher and its helpers
func (s *Searcher) TotalNodes() int {
	total := s.Nodes
	for _, h := range s.helpers {
		total += int(h.published.Load())
	}
	return total
}

// sets number of search threads (this searcher counts as one)
func (s *Searcher) SetThreads(n int) {
	if s.NN == nil {
		return
	}

	want := min(max(n, 1), MAX_THREADS) - 1

	if want < len(s.helpers) {
		clear(s.helpers[want:]) // let removed searchers be collected
		s.helpers = s.helpers[:want]
	}

	for len(s.helpers) < want {
		nn := &nnue.NNUE{Network: s.NN.Network}
		s.helpers = append(s.helpers, NewSearcher(nn, s.TT, s.Stop))
	}
}

// number of search threads, including this one
func (s *Searcher) Threads() int {
	return len(s.helpers) + 1
}

// starts the helper threads on a copy of the position b and returns at once
func (s *Searcher) StartHelpers(b *board.Board, nodeLimit int, wg *sync.WaitGroup) {
	for i, h := range s.helpers {
		// sync with the main thread's current state
		h.TT = s.TT
		h.Stop = s.Stop
		h.NodeLimit = nodeLimit
		h.lead = s
		h.Nodes = 0
		h.published.Store(0)
		h.SetGame(s.game)

		root := *b
		id := i + 1

		wg.Go(func() {
			h.runHelper(root, id)
		})

	}
}

// iterative deepening for a helper thread
func (s *Searcher) runHelper(root board.Board, id int) {
	b := root
	s.NN.Reset(&b)
	s.Reset()

	score := 0

	for d := 1; d <= helperMaxDepth; d++ {
		if skipDepth(id, d) {
			continue
		}

		_, sc := s.Search(&b, d, score)

		if s.Stop.Load() {
			return
		}

		score = sc
	}
}
