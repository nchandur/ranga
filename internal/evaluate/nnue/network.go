package nnue

type Network struct {
	FeatureWeights [768]Accumulator
	FeatureBias    Accumulator
	OutputWeights  [2 * HiddenSize]int16
	OutputBias     int16
}

func (net *Network) Evaluate(us, them *Accumulator) int32 {
	output := int32(0)

	// side to move accumulator
	for i := range HiddenSize {
		output += screlu(us.Values[i]) * int32(net.OutputWeights[i])
	}

	// not side to move accumulator
	for i := range HiddenSize {
		output += screlu(them.Values[i]) * int32(net.OutputWeights[i])
	}

	// reduce quantization
	output /= QA

	// add bias
	output += int32(net.OutputBias)

	// apply eval scale
	output *= Scale

	// remove quantization
	output /= (QA * QB)

	return output
}
