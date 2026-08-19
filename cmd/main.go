package main

import (
	"context"
	"fmt"
	"log"
	"ranga/internal/board"
	"ranga/internal/evaluate"
	"ranga/internal/search"
	"time"
)

var version string = ""

func main() {
	b := board.NewBoard()

	if err := b.ParseFEN("rnb1kbnr/ppp1pppp/8/4q3/8/2N5/PPPP1PPP/R1BQKBNR w KQkq - 2 4"); err != nil {
		log.Fatal(err)
	}

	searcher := search.NewSearcher(evaluate.Evaluator{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*time.Second))
	defer cancel()
	move := searcher.Search(ctx, &b, 3)

	fmt.Printf("%s\n", move)

}
