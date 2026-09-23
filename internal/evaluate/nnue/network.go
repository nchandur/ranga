package nnue

import (
	"math/rand"
)

type Network struct {
	FeatureWeights [768]Accumulator
	FeatureBias    Accumulator
	OutputWeights  [2 * HiddenSize]int16
	OutputBias     int16
}

// fills the network with small random weights: purely for wiring up search
func (net *Network) Randomize() {
	for f := range net.FeatureWeights {
		for i := range HiddenSize {
			net.FeatureWeights[f].Values[i] = int16(rand.Intn(201) - 100)
		}
	}
	for i := range HiddenSize {
		net.FeatureBias.Values[i] = int16(rand.Intn(201) - 100)
	}
	for i := range net.OutputWeights {
		net.OutputWeights[i] = int16(rand.Intn(201) - 100)
	}
	net.OutputBias = int16(rand.Intn(201) - 100)
}

func (net *Network) Evaluate(us, them *Accumulator) int32 {
	output := int32(0)

	// side to move accumulator
	for i := range HiddenSize {
		output += screlu(us.Values[i]) * int32(net.OutputWeights[i])
	}

	// not side to move accumulator
	for i := range HiddenSize {
		output += screlu(them.Values[i]) * int32(net.OutputWeights[HiddenSize+i])
	}

	// reduce quantization
	output /= int32(QA)

	// add bias
	output += int32(net.OutputBias)

	// apply eval scale
	output *= Scale

	// remove quantization
	output /= (int32(QA) * int32(QB))

	return output
}
