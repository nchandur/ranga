package nnue

import (
	"math/rand"
)

type Network struct {
	FeatureWeights [768]Accumulator
	FeatureBias    Accumulator
	OutputWeights  [2 * HiddenSize]int
	OutputBias     int
}

// fills the network with small random weights: purely for wiring up search
func (net *Network) Randomize() {
	for f := range net.FeatureWeights {
		for i := range HiddenSize {
			net.FeatureWeights[f].Values[i] = rand.Intn(201) - 100
		}
	}
	for i := range HiddenSize {
		net.FeatureBias.Values[i] = rand.Intn(201) - 100
	}
	for i := range net.OutputWeights {
		net.OutputWeights[i] = rand.Intn(201) - 100
	}
	net.OutputBias = rand.Intn(201) - 100
}

func (net *Network) Evaluate(us, them *Accumulator) int {
	output := 0

	// side to move accumulator
	for i := range HiddenSize {
		output += screlu(us.Values[i]) * net.OutputWeights[i]
	}

	// not side to move accumulator
	for i := range HiddenSize {
		output += screlu(them.Values[i]) * net.OutputWeights[HiddenSize+i]
	}

	// reduce quantization
	output /= QA

	// add bias
	output += net.OutputBias

	// apply eval scale
	output *= Scale

	// remove quantization
	output /= (QA * QB)

	return output
}
