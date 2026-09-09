package main

import (
	"fmt"
	"log"
	"ranga/internal/pgn"
)

var version string = ""

func main() {

	gamePGN := `
	[Event "Fastchess Tournament"]
	[Site "?"]
	[Date "2026.09.03"]
	[Round "5"]
	[White "ranga-1.12-a"]
	[Black "ranga-1.12-b"]
	[Result "0-1"]
	[SetUp "1"]
	[FEN "r3kbnr/1bqppppp/p1p5/8/4P3/1P1B4/P1P2PPP/RNBQK2R w KQkq - 1 8"]
	[GameDuration "00:00:00"]
	[GameStartTime "2026-09-03T23:48:44 +0000"]
	[GameEndTime "2026-09-03T23:48:45 +0000"]
	[PlyCount "36"]
	[Termination "normal"]
	[TimeControl "-"]

	8. O-O {+0.43/4 0.026s} e5 {-0.77/4 0.025s} 9. Be3 {+0.35/4 0.025s}
	d5 {-0.47/4 0.026s} 10. c3 {+0.09/4 0.023s} O-O-O {-0.21/4 0.021s}
	11. exd5 {+0.26/4 0.023s} cxd5 {+0.09/5 0.024s} 12. Bf5+ {-0.03/4 0.025s}
	Kb8 {+0.37/5 0.028s} 13. h3 {-0.11/3 0.027s} Ne7 {+0.45/4 0.023s}
	14. Bg4 {-0.74/4 0.024s} Ng6 {+0.56/4 0.025s} 15. Bf5 {-0.56/3 0.028s}
	Nh4 {+0.70/4 0.029s} 16. Qh5 {-0.99/4 0.022s} Be7 {+1.03/4 0.025s}
	17. Bxh7 {-1.03/3 0.025s} g6 {+3.58/4 0.033s} 18. Bb6 {-3.66/4 0.019s}
	Qxb6 {+3.90/5 0.019s} 19. Qxe5+ {-3.66/5 0.027s} Bd6 {+3.66/5 0.025s}
	20. Qg7 {-3.60/5 0.024s} d4 {+3.60/3 0.027s} 21. cxd4 {-3.90/3 0.036s}
	Bxg2 {+3.90/3 0.020s} 22. Rd1 {-3.72/3 0.023s} Qc7 {+5.10/3 0.026s}
	23. f4 {-5.94/4 0.027s} Bxf4 {+5.94/3 0.033s} 24. Rd3 {-5.94/3 0.033s}
	Qc1+ {+9999.97/4 0.025s} 25. Kf2 {-9999.98/7 0.025s}
	Qf1# {+9999.99/4 0.022s, Black mates} 0-1
	`

	game := pgn.NewGame(gamePGN)

	game.Board.Print()

	lines, err := game.GenerateLines()
	if err != nil {
		log.Fatal(err)
	}
	for idx, l := range lines {
		fmt.Println(idx+1, " ", l)
	}

	// engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	// engine.Run()

}
