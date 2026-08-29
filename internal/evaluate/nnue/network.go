package nnue

import (
	"math/rand/v2"
	"ranga/internal/board"
)

// holds the static weights and biases for the (768->256)x2->1 NNUE
type Network struct {
	FeatureWeights [InputSize][HiddenSize]int16
	FeatureBias    [HiddenSize]int16
	OutputWeights  [HiddenSize * 2]int16
	OutputBias     int32
}

// generates network with randomized weights (FOR TESTING ONLY. REMOVE THIS SHIT AFTER TRAINING)
func NewRandomNetwork() *Network {
	nn := &Network{}

	// populate feature layer with random integers
	for i := range InputSize {
		for j := range HiddenSize {
			nn.FeatureWeights[i][j] = int16(rand.Int64N(200) - 100)
		}
	}

	for j := range HiddenSize {
		nn.FeatureBias[j] = int16(rand.Int64N(200) - 100)
	}

	// populate outputl layer
	for i := range HiddenSize * 2 {
		nn.OutputWeights[i] = int16(rand.Int64N(200) - 100)
	}

	nn.OutputBias = int32(rand.Int64N(200) - 100)

	return nn
}

// computes the final network score from the accumulator.
func (nn *Network) Evaluate(acc *Accumulator, sideToMove board.Color) int {
	var output int32 = nn.OutputBias

	// determine perspective
	var us, them *[HiddenSize]int16
	if sideToMove == board.White {
		us = &acc.White
		them = &acc.Black
	} else {
		us = &acc.Black
		them = &acc.White
	}

	// activated 'us' features by first half of output weights
	for i := range HiddenSize {
		activated := int32(crelu(us[i]))
		output += activated * int32(nn.OutputWeights[i])
	}

	// activated 'them' features by second half of output weights
	for i := range HiddenSize {
		activated := int32(crelu(them[i]))
		output += activated * int32(nn.OutputWeights[HiddenSize+i])
	}

	// scale  down to roughly match standard centipawn scaling (TWEAK DIVISOR AFTER TRAINING)
	return int(output / 4096)
}
