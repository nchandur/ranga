package nnue

import "testing"

func newMockNetwork() *Network {
	net := &Network{}

	for i := range net.FeatureBias.Values {
		net.FeatureBias.Values[i] = int16(i % 100)
	}

	for f := range 5 {
		for i := range net.FeatureWeights[f].Values {
			net.FeatureWeights[f].Values[i] = int16((f + 1) * (i + 1))
		}
	}

	return net
}

func TestAccumulator(t *testing.T) {
	net := newMockNetwork()

	t.Run("NewAccumulator initializes with FeatureBias", func(t *testing.T) {
		acc := NewAccumulator(net)

		for i := range HiddenSize {
			if acc.Values[i] != net.FeatureBias.Values[i] {
				t.Fatalf("acc.Values[%d] = %d; want bias %d",
					i, acc.Values[i], net.FeatureBias.Values[i])
			}
		}
	})

	t.Run("AddFeature adds weights to accumulator values", func(t *testing.T) {
		acc := NewAccumulator(net)
		featureIdx := 0

		acc.AddFeature(net, featureIdx)

		for i := range HiddenSize {
			want := net.FeatureBias.Values[i] + net.FeatureWeights[featureIdx].Values[i]
			if acc.Values[i] != want {
				t.Fatalf("AddFeature index %d mismatch: got %d, want %d",
					i, acc.Values[i], want)
			}
		}
	})

	t.Run("AddSub processes 1 addition and 1 subtraction", func(t *testing.T) {
		acc := NewAccumulator(net)
		acc.AddSub(net, 1, 0)

		for i := range HiddenSize {
			want := net.FeatureBias.Values[i] + net.FeatureWeights[1].Values[i] - net.FeatureWeights[0].Values[i]
			if acc.Values[i] != want {
				t.Fatalf("AddSub index %d mismatch: got %d, want %d",
					i, acc.Values[i], want)
			}
		}
	})

	t.Run("AddSub2 processes 1 addition and 2 subtractions", func(t *testing.T) {
		acc := NewAccumulator(net)
		acc.AddSub2(net, 2, 1, 0)

		for i := range HiddenSize {
			want := net.FeatureBias.Values[i] + net.FeatureWeights[2].Values[i] - net.FeatureWeights[1].Values[i] - net.FeatureWeights[0].Values[i]
			if acc.Values[i] != want {
				t.Fatalf("AddSub2 index %d mismatch: got %d, want %d",
					i, acc.Values[i], want)
			}
		}
	})

	t.Run("Add2Sub2 processes 2 additions and 2 subtractions", func(t *testing.T) {
		acc := NewAccumulator(net)
		acc.Add2Sub2(net, 3, 4, 1, 0)

		for i := range HiddenSize {
			want := net.FeatureBias.Values[i] + net.FeatureWeights[3].Values[i] + net.FeatureWeights[4].Values[i] - net.FeatureWeights[1].Values[i] - net.FeatureWeights[0].Values[i]
			if acc.Values[i] != want {
				t.Fatalf("Add2Sub2 index %d mismatch: got %d, want %d",
					i, acc.Values[i], want)
			}
		}
	})
}
