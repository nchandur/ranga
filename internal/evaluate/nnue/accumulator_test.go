package nnue

import "testing"

func newMockNetwork() *Network {
	net := &Network{}

	for i := range net.FeatureBias.Values {
		net.FeatureBias.Values[i] = int16(i % 100)
	}

	for f := 0; f < 3; f++ {
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

	t.Run("RemoveFeature subtracts weights from accumulator values", func(t *testing.T) {
		acc := NewAccumulator(net)
		featureIdx := 1

		acc.RemoveFeature(net, featureIdx)

		for i := range HiddenSize {
			want := net.FeatureBias.Values[i] - net.FeatureWeights[featureIdx].Values[i]
			if acc.Values[i] != want {
				t.Fatalf("RemoveFeature index %d mismatch: got %d, want %d",
					i, acc.Values[i], want)
			}
		}
	})

	t.Run("AddFeature followed by RemoveFeature restores original state", func(t *testing.T) {
		acc := NewAccumulator(net)
		initial := acc

		acc.AddFeature(net, 0)
		acc.AddFeature(net, 1)
		acc.RemoveFeature(net, 1)
		acc.RemoveFeature(net, 0)

		for i := range HiddenSize {
			if acc.Values[i] != initial.Values[i] {
				t.Fatalf("invariant broken at index %d: got %d, want initial %d",
					i, acc.Values[i], initial.Values[i])
			}
		}
	})

	t.Run("Multiple feature accumulations match sum", func(t *testing.T) {
		acc := NewAccumulator(net)

		acc.AddFeature(net, 0)
		acc.AddFeature(net, 1)

		for i := range HiddenSize {
			want := net.FeatureBias.Values[i] +
				net.FeatureWeights[0].Values[i] +
				net.FeatureWeights[1].Values[i]

			if acc.Values[i] != want {
				t.Fatalf("cumulative sum mismatch at index %d: got %d, want %d",
					i, acc.Values[i], want)
			}
		}
	})
}
