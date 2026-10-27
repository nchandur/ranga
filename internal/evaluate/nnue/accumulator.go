package nnue

type Accumulator struct {
	Values [HiddenSize]int16
}

// initialized with bias for efficient operations later
func NewAccumulator(net *Network) Accumulator {
	return net.FeatureBias
}

// add feature to accumulator
func (acc *Accumulator) AddFeature(net *Network, featureIdx int) {
	for i := range acc.Values {
		d := net.FeatureWeights[featureIdx].Values[i]
		acc.Values[i] += d
	}
}

// handles 1 addition and 1 subtraction (quiet moves)
func (acc *Accumulator) AddSub(net *Network, addIdx, subIdx int) {
	for i := range acc.Values {
		acc.Values[i] += net.FeatureWeights[addIdx].Values[i] - net.FeatureWeights[subIdx].Values[i]
	}
}

// handles 1 addition and 2 subtractions (captures and enpassant)
func (acc *Accumulator) AddSub2(net *Network, addIdx, sub1Idx, sub2Idx int) {
	for i := range acc.Values {
		acc.Values[i] += net.FeatureWeights[addIdx].Values[i] - net.FeatureWeights[sub1Idx].Values[i] - net.FeatureWeights[sub2Idx].Values[i]
	}
}

// handles 2 additions and 2 subtractions (castling)
func (acc *Accumulator) Add2Sub2(net *Network, add1Idx, add2Idx, sub1Idx, sub2Idx int) {
	for i := range acc.Values {
		acc.Values[i] += net.FeatureWeights[add1Idx].Values[i] + net.FeatureWeights[add2Idx].Values[i] - net.FeatureWeights[sub1Idx].Values[i] - net.FeatureWeights[sub2Idx].Values[i]
	}
}
