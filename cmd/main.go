package main

import (
	"os"
	"ranga/internal/uci"
)

var version string = ""

func main() {

	// nn, err := nnue.LoadEmbeddedNetwork()
	// fmt.Println("load err:", err)
	// fmt.Println("sample weight:", nn.FeatureWeights[0][0], nn.FeatureWeights[100][50])
	// fmt.Println("output bias:", nn.OutputBias)

	// b := board.NewBoard()
	// err = b.ParseFEN("r1bqkbnr/pp1ppp1p/n1p3p1/8/3P2P1/8/PPPQPP1P/RNB1KBNR w KQkq - 0 1")

	// acc := &nnue.Accumulator{}
	// acc.Refresh(nn, &b)
	// fmt.Println("acc.White[0..3]:", acc.White[0], acc.White[1], acc.White[2], acc.White[3])
	// fmt.Println("eval:", nn.Evaluate(acc, board.White))

	engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	engine.Run()

}
