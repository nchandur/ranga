package nnue

import (
	"testing"
)

func TestNetwork(t *testing.T) {
	t.Run("Evaluate with zero accumulators evaluates OutputBias and Scale", func(t *testing.T) {
		for bucket := range OutputBuckets {
			net := &Network{}
			net.OutputBias[bucket] = 50
			us := &Accumulator{}
			them := &Accumulator{}
			got := net.Evaluate(us, them, bucket)
			want := (int32(net.OutputBias[bucket]) * Scale) / (int32(QA) * int32(QB))

			if got != want {
				t.Errorf("bucket %d: Evaluate with zero accumulators = %d; want %d", bucket, got, want)
			}
		}
	})

	t.Run("Evaluate separates us and them halves accurately", func(t *testing.T) {
		for bucket := range OutputBuckets {
			net := &Network{}
			net.OutputBias[bucket] = 0
			net.OutputWeights[bucket][0] = 10
			net.OutputWeights[bucket][HiddenSize] = 20

			us := &Accumulator{}
			them := &Accumulator{}

			us.Values[0] = 4
			them.Values[0] = 5

			sum := int32(16*10 + 25*20)
			want := ((sum / int32(QA)) * Scale) / (int32(QA) * int32(QB))

			got := net.Evaluate(us, them, bucket)
			if got != want {
				t.Errorf("bucket %d: Evaluate mismatch:\ngot:  %d\nwant: %d", bucket, got, want)
			}
		}
	})

	t.Run("Evaluate is sensitive to perspective order", func(t *testing.T) {
		for bucket := range OutputBuckets {
			net := &Network{}
			net.OutputBias[bucket] = 0

			net.OutputWeights[bucket][0] = 100
			net.OutputWeights[bucket][HiddenSize] = -100

			accA := &Accumulator{}
			accB := &Accumulator{}

			accA.Values[0] = QA
			accB.Values[0] = 0

			evalA := net.Evaluate(accA, accB, bucket)
			evalB := net.Evaluate(accB, accA, bucket)

			if evalA == evalB {
				t.Errorf("bucket %d: expected different evaluations when swapping us and them perspectives; both got %d", bucket, evalA)
			}
			if evalA <= 0 {
				t.Errorf("bucket %d: expected positive evaluation for evalA, got %d", bucket, evalA)
			}
		}
	})

	t.Run("Evaluate uses only the selected bucket's bias", func(t *testing.T) {
		net := &Network{}
		for b := range OutputBuckets {
			net.OutputBias[b] = int16((b + 1) * 2000)
		}
		us := &Accumulator{}
		them := &Accumulator{}

		for bucket := range OutputBuckets {
			got := net.Evaluate(us, them, bucket)
			want := (int32(net.OutputBias[bucket]) * Scale) / (int32(QA) * int32(QB))
			if got != want {
				t.Errorf("bucket %d: Evaluate = %d; want %d", bucket, got, want)
			}
		}
	})

	t.Run("Evaluate uses only the selected bucket's weights", func(t *testing.T) {
		net := &Network{}
		// only bucket 3 has non-zero output weights
		const active = 3
		net.OutputWeights[active][0] = 100
		net.OutputWeights[active][HiddenSize] = 100

		us := &Accumulator{}
		them := &Accumulator{}
		us.Values[0] = QA
		them.Values[0] = QA

		for bucket := range OutputBuckets {
			got := net.Evaluate(us, them, bucket)
			if bucket == active {
				if got <= 0 {
					t.Errorf("bucket %d: expected positive evaluation, got %d", bucket, got)
				}
			} else if got != 0 {
				t.Errorf("bucket %d: expected 0 from empty bucket, got %d", bucket, got)
			}
		}
	})
}
