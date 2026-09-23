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

// remove feature from accumulator
func (acc *Accumulator) RemoveFeature(net *Network, featureIdx int) {
	for i := range acc.Values {
		d := net.FeatureWeights[featureIdx].Values[i]
		acc.Values[i] -= d
	}
}
