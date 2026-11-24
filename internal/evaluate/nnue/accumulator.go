package nnue

type Accumulator struct {
	Values [HiddenSize]int16
}

// initialized with bias for efficient operations later
func NewAccumulator(net *Network) Accumulator {
	return net.FeatureBias
}

func (acc *Accumulator) AddFeature(net *Network, featureIdx int) {
	a := &net.FeatureWeights[featureIdx].Values
	for i := range acc.Values {
		acc.Values[i] += a[i]
	}
}

// dst = src + add - sub (quiet moves)
func (dst *Accumulator) AddSubFrom(src *Accumulator, net *Network, addIdx, subIdx int) {
	a := &net.FeatureWeights[addIdx].Values
	s := &net.FeatureWeights[subIdx].Values
	for i := range dst.Values {
		dst.Values[i] = src.Values[i] + a[i] - s[i]
	}
}

// dst = src + add - sub1 - sub2 (captures, en passant)
func (dst *Accumulator) AddSub2From(src *Accumulator, net *Network, addIdx, sub1Idx, sub2Idx int) {
	a := &net.FeatureWeights[addIdx].Values
	s1 := &net.FeatureWeights[sub1Idx].Values
	s2 := &net.FeatureWeights[sub2Idx].Values
	for i := range dst.Values {
		dst.Values[i] = src.Values[i] + a[i] - s1[i] - s2[i]
	}
}

// dst = src + add1 + add2 - sub1 - sub2 (castling)
func (dst *Accumulator) Add2Sub2From(src *Accumulator, net *Network, add1Idx, add2Idx, sub1Idx, sub2Idx int) {
	a1 := &net.FeatureWeights[add1Idx].Values
	a2 := &net.FeatureWeights[add2Idx].Values
	s1 := &net.FeatureWeights[sub1Idx].Values
	s2 := &net.FeatureWeights[sub2Idx].Values
	for i := range dst.Values {
		dst.Values[i] = src.Values[i] + a1[i] + a2[i] - s1[i] - s2[i]
	}
}

// in-place variants
func (acc *Accumulator) AddSub(net *Network, addIdx, subIdx int) {
	acc.AddSubFrom(acc, net, addIdx, subIdx)
}

func (acc *Accumulator) AddSub2(net *Network, addIdx, sub1Idx, sub2Idx int) {
	acc.AddSub2From(acc, net, addIdx, sub1Idx, sub2Idx)
}

func (acc *Accumulator) Add2Sub2(net *Network, add1Idx, add2Idx, sub1Idx, sub2Idx int) {
	acc.Add2Sub2From(acc, net, add1Idx, add2Idx, sub1Idx, sub2Idx)
}
