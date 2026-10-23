package nnue

import (
	"math/rand"
)

type Network struct {
	FeatureWeights [768]Accumulator
	FeatureBias    Accumulator
	OutputWeights  [OutputBuckets][2 * HiddenSize]int16
	OutputBias     [OutputBuckets]int16
}

// fills the network with small random weights: purely for wiring up search
func (net *Network) Randomize() {
	for b := range OutputBuckets {
		for i := range net.OutputWeights[b] {
			net.OutputWeights[b][i] = int16(rand.Intn(201) - 100)
		}
		net.OutputBias[b] = int16(rand.Intn(201) - 100)
	}
}

func (net *Network) Evaluate(us, them *Accumulator, bucket int) int32 {
	w := &net.OutputWeights[bucket]
	output := lizardDot(us.Values[:], w[:HiddenSize])
	output += lizardDot(them.Values[:], w[HiddenSize:])

	// reduce quantization
	output /= int32(QA)

	// add bias
	output += int32(net.OutputBias[bucket])

	// apply eval scale
	output *= Scale

	// remove quantization
	output /= (int32(QA) * int32(QB))

	return output
}
