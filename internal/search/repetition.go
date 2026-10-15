package search

import "ranga/internal/board"

// forgets the game history
func (s *Searcher) ClearGame() {
	s.game = s.game[:0]
}

// records a move played in the game from the position with key prevKey
func (s *Searcher) AddGameMove(prevKey uint64, fiftyMove int) {
	// irreversible move: nothing before it can ever repeat
	if fiftyMove == 0 {
		s.game = s.game[:0]
		return
	}

	s.game = append(s.game, prevKey)
}

// replaces the game history with copy of keys
func (s *Searcher) SetGame(keys []uint64) {
	s.game = append(s.game[:0], keys...)
}

// checks whether current board position has occurred previously
func (s *Searcher) IsRepetition(b *board.Board) bool {
	games := len(s.game)
	cur := games + b.Ply

	for i := max(cur-b.FiftyMove, 0); i < cur; i++ {
		if i < games {
			if s.game[i] == b.Key {
				return true
			}
		} else if s.path[i-games] == b.Key {
			return true
		}
	}

	return false
}
