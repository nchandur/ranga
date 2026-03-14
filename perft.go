package main

import (
	"fmt"
	"time"
)

var perftLeafNodes int64

func (b *Board) Perft(depth int) {

	if depth == 0 {
		perftLeafNodes++
		return
	}

	b.GenMoves()

	for idx := b.MoveListStart[b.Ply]; idx < b.MoveListStart[b.Ply+1]; idx++ {

		move := Move(b.MoveList[idx])

		if !b.MakeMove(move) {
			continue
		}

		b.Perft(depth - 1)

		b.TakeMove()
	}
}

func (b *Board) RunPerft(depth int) {

	perftLeafNodes = 0

	start := time.Now()

	b.Perft(depth)

	elapsed := time.Since(start)

	nps := int64(float64(perftLeafNodes) / elapsed.Seconds())

	fmt.Printf("\nPerft Test Complete\n")
	fmt.Printf("-------------------\n")
	fmt.Printf("Depth: %d\n", depth)
	fmt.Printf("Nodes: %d\n", perftLeafNodes)
	fmt.Printf("Time:  %s\n", elapsed)
	fmt.Printf("NPS:   %d\n\n", nps)
}


func (b *Board) PerftDivide(depth int) {

	perftLeafNodes = 0

	start := time.Now()

	b.GenMoves()

	moveNum := 0

	for idx := b.MoveListStart[b.Ply]; idx < b.MoveListStart[b.Ply+1]; idx++ {

		move := Move(b.MoveList[idx])

		if !b.MakeMove(move) {
			continue
		}

		moveNum++

		nodesBefore := perftLeafNodes

		b.Perft(depth - 1)

		b.TakeMove()

		nodes := perftLeafNodes - nodesBefore

		fmt.Printf("%2d: %-6s %d\n", moveNum, move.String(), nodes)
	}

	elapsed := time.Since(start)

	nps := int64(float64(perftLeafNodes) / elapsed.Seconds())

	fmt.Printf("\nPerft Divide Complete\n")
	fmt.Printf("---------------------\n")
	fmt.Printf("Depth: %d\n", depth)
	fmt.Printf("Nodes: %d\n", perftLeafNodes)
	fmt.Printf("Time:  %s\n", elapsed)
	fmt.Printf("NPS:   %d\n\n", nps)
}