package board

import (
	"context"
	"fmt"
)

func Perft(ctx context.Context, b *Board, depth int) int64 {

	if ctx.Err() != nil {
		return 0
	}

	if depth == 0 {
		return 1
	}

	var nodes int64 = 0
	var list MoveList

	list.GenerateMoves(b)

	for i := range list.Count {
		move := list.Moves[i]

		copy := b.Preserve()

		if !b.MakeMove(move, false) {
			b.Restore(&copy)
			continue
		}

		count := Perft(ctx, b, depth-1)
		nodes += count

		b.Restore(&copy)

		if depth > 2 && ctx.Err() != nil {
			return 0
		}

	}

	return nodes
}

func PerftDivide(b *Board, depth int) int64 {
	if depth == 0 {
		return 1
	}

	ctx := context.Background()

	var totalNodes int64 = 0
	var list MoveList
	list.GenerateMoves(b)

	for i := range list.Count {
		move := list.Moves[i]
		copy := b.Preserve()

		if !b.MakeMove(move, false) {
			b.Restore(&copy)
			continue
		}

		nodes := Perft(ctx, b, depth-1)
		totalNodes += nodes

		b.Restore(&copy)

		fmt.Println(move, ":", nodes)
	}

	fmt.Println("\nTotal nodes:", totalNodes)
	return totalNodes
}
