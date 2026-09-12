package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"ranga/internal/board"
	"regexp"
	"strings"
)

var whiteRe = regexp.MustCompile(`\[White\s*"(.+)"\]`)
var blackRe = regexp.MustCompile(`\[Black\s*"(.+)"\]`)
var resultRe = regexp.MustCompile(`\[Result\s*"(.+)"\]`)
var fenRe = regexp.MustCompile(`\[FEN\s*"(.+)"\]`)
var taglineRe = regexp.MustCompile(`(?m)^\s*\[.*\]\s*$`)

// marks the start of a new game header block.
var eventTagPattern = regexp.MustCompile(`(?m)^\[Event\s`)

type Game struct {
	White  string
	Black  string
	Result string
	FEN    string // starting fen position
	Moves  []string
	Scores []int
	board.Board
}

func NewGame(pgn string) (Game, error) {

	var game Game
	game.Board = board.NewBoard()
	game.parsePlayers(pgn)
	game.parseResult(pgn)
	game.parseFEN(pgn)

	if err := game.Board.ParseFEN(game.FEN); err != nil {
		return Game{}, err
	}

	game.parseMoves(pgn)

	return game, nil
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
		switch matches[1] {
		case "1-0":
			g.Result = "1.0"
		case "0-1":
			g.Result = "0.0"
		case "1/2-1/2":
			g.Result = "0.5"
		}
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
	g.Scores = make([]int, len(pairs))
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
		lines = append(lines, fmt.Sprintf("%s | %d | %s", g.Board.FEN(), g.Scores[i], g.Result))
	}

	return lines, nil
}

// splits a file containing multiple concatenated PGN games into individual games
func SplitGames(data string) []string {
	idx := eventTagPattern.FindAllStringIndex(data, -1)
	if len(idx) == 0 {
		return nil
	}

	games := make([]string, 0, len(idx))
	for i, m := range idx {
		start := m[0]
		end := len(data)
		if i+1 < len(idx) {
			end = idx[i+1][0]
		}
		game := strings.TrimSpace(data[start:end])
		if game != "" {
			games = append(games, game)
		}
	}
	return games
}

// reads a file containing n concatenated PGN games and returns the generated "<FEN> | <score> | <result>"
func ParseGamesFromFile(source string) ([]string, []error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return nil, []error{fmt.Errorf("reading file: %w", err)}
	}

	pgnGames := SplitGames(string(data))
	totalGames := len(pgnGames)

	var allLines []string
	var errs []error
	for i, pgnStr := range pgnGames {
		game, err := NewGame(pgnStr)
		if err != nil {
			errs = append(errs, fmt.Errorf("game %d: %w", i+1, err))
		}

		lines, err := game.GenerateLines()
		if err != nil {
			errs = append(errs, fmt.Errorf("game %d: %w", i+1, err))
			continue
		}
		allLines = append(allLines, lines...)
		fmt.Printf("\r\033[K[game %d/%d] positions parsed: %d", i+1, totalGames, len(allLines))
		os.Stdout.Sync()
	}

	return allLines, errs
}

// write FEN strings to file
func WriteGamesToFile(destination string, games []string) error {
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, game := range games {
		if _, err := w.WriteString(game); err != nil {
			return fmt.Errorf("writing line: %w", err)
		}
		if err := w.WriteByte('\n'); err != nil {
			return fmt.Errorf("writing newline: %w", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flushing output file: %w", err)
	}

	return nil
}

// write errors to file
func WriteErrsToFile(destination string, errors []error) error {
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range errors {

		if e != nil {
			if _, err := w.WriteString(e.Error()); err != nil {
				return fmt.Errorf("writing line: %w", err)
			}
			if err := w.WriteByte('\n'); err != nil {
				return fmt.Errorf("writing newline: %w", err)
			}
		}

	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flushing output file: %w", err)
	}

	return nil
}
