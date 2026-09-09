package pgn

import (
	"fmt"
	"log"
	"ranga/internal/board"
	"regexp"
)

var whiteRe = regexp.MustCompile(`\[White\s*"(.+)"\]`)
var blackRe = regexp.MustCompile(`\[Black\s*"(.+)"\]`)
var resultRe = regexp.MustCompile(`\[Result\s*"(.+)"\]`)
var fenRe = regexp.MustCompile(`\[FEN\s*"(.+)"\]`)
var taglineRe = regexp.MustCompile(`(?m)^\s*\[.*\]\s*$`)

type Game struct {
	White  string
	Black  string
	Result string
	FEN    string // starting fen position
	Moves  []string
	Scores []string
	board.Board
}

func NewGame(pgn string) Game {

	var game Game
	game.Board = board.NewBoard()
	game.parsePlayers(pgn)
	game.parseResult(pgn)
	game.parseFEN(pgn)

	if err := game.Board.ParseFEN(game.FEN); err != nil {
		log.Fatal(err)
	}

	game.parseMoves(pgn)

	return game
}

// parse players
func (g *Game) parsePlayers(pgn string) {
	if matches := whiteRe.FindStringSubmatch(pgn); len(matches) > 0 {
		g.White = matches[1]
	}

	if matches := blackRe.FindStringSubmatch(pgn); len(matches) > 0 {
		g.Black = matches[1]
	}
}

// parse result
func (g *Game) parseResult(pgn string) {
	if matches := resultRe.FindStringSubmatch(pgn); len(matches) > 0 {
		g.Result = matches[1]
	}
}

// parse starting fen
func (g *Game) parseFEN(pgn string) {
	if matches := fenRe.FindStringSubmatch(pgn); len(matches) > 0 {
		g.FEN = matches[1]
	}
}

// parse moves
func (g *Game) parseMoves(pgn string) {
	movetext := extractMovetext(pgn)

	pairs, err := tokenizeWithComments(movetext)
	if err != nil {
		log.Fatalf("failed to tokenize movetext: %v", err)
	}

	g.Moves = make([]string, len(pairs))
	g.Scores = make([]string, len(pairs))
	for i, p := range pairs {
		g.Moves[i] = p.san
		g.Scores[i] = extractScore(p.comment)
	}
}

// returns everything in the PGN after the last header tag
func extractMovetext(pgn string) string {
	matches := taglineRe.FindAllStringIndex(pgn, -1)
	if len(matches) == 0 {
		return pgn
	}
	last := matches[len(matches)-1]
	return pgn[last[1]:]
}

// replays moves on board returns one line per ply
func (g *Game) GenerateLines() ([]string, error) {
	lines := make([]string, 0, len(g.Moves))

	for i, san := range g.Moves {
		m := ParseSAN(&g.Board, san)
		if m == board.NOMOVE {
			return nil, fmt.Errorf("move %d (%q): unparseable in position\n%s", i+1, san, g.Board.FEN())
		}
		if ok := g.Board.MakeMove(m, false); !ok {
			return nil, fmt.Errorf("move %d (%q): illegal move\n%s", i+1, san, g.Board.FEN())
		}
		g.Board.Ply++
		lines = append(lines, fmt.Sprintf("%s | %s | %s", g.Board.FEN(), g.Scores[i], g.Result))
	}

	return lines, nil
}
