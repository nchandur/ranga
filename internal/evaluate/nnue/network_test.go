package nnue

import (
	"testing"
)

func TestNetwork(t *testing.T) {
	t.Run("Evaluate with zero accumulators evaluates OutputBias and Scale", func(t *testing.T) {
		net := &Network{}
		net.OutputBias = 50
		us := &Accumulator{}
		them := &Accumulator{}
		got := net.Evaluate(us, them)
		want := (int32(net.OutputBias) * Scale) / (int32(QA) * int32(QB))

		if got != want {
			t.Errorf("Evaluate with zero accumulators = %d; want %d", got, want)
		}
	})

	t.Run("Evaluate separates us and them halves accurately", func(t *testing.T) {
		net := &Network{}
		net.OutputBias = 0
		net.OutputWeights[0] = 10
		net.OutputWeights[HiddenSize] = 20

		us := &Accumulator{}
		them := &Accumulator{}

		us.Values[0] = 4
		them.Values[0] = 5

		sum := int32(16*10 + 25*20)
		want := ((sum / int32(QA)) * Scale) / (int32(QA) * int32(QB))

		got := net.Evaluate(us, them)
		if got != want {
			t.Errorf("Evaluate mismatch:\ngot:  %d\nwant: %d", got, want)
		}
	})
	t.Run("Evaluate is sensitive to perspective order", func(t *testing.T) {
		net := &Network{}
		net.OutputBias = 0

		net.OutputWeights[0] = 100
		net.OutputWeights[HiddenSize] = -100

		accA := &Accumulator{}
		accB := &Accumulator{}

		accA.Values[0] = QA
		accB.Values[0] = 0

		evalA := net.Evaluate(accA, accB)
		evalB := net.Evaluate(accB, accA)

		if evalA == evalB {
			t.Errorf("expected different evaluations when swapping us and them perspectives; both got %d", evalA)
		}
		if evalA <= 0 {
			t.Errorf("expected positive evaluation for evalA, got %d", evalA)
		}
	})
}

func TestOutputWeightsFitLizard(t *testing.T) {
	net, _ := LoadEmbedded()

	var maxAbs int32
	for _, w := range net.OutputWeights {
		a := int32(w)
		if a < 0 {
			a = -a
		}
		if a > maxAbs {
			maxAbs = a
		}
	}

	limit := int32(32767) / int32(QA)
	t.Logf("max |output weight| = %d, limit = %d", maxAbs, limit)
	if maxAbs > limit {
		t.Fatalf("net not compatible with Lizard SCReLU")
	}
}
