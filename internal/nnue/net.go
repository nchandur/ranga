package nnue

import "ranga/internal/board"

// hidden layer sizes for two small fully-connected layers after feature transformer
const (
	L2Size = 32
	L3Size = 32
)

// holds every weight matrix
type Network struct {
	FTWeights []float32 // [InputSize][HiddenSize], row-major: feature i -> FTWeights[i*HiddenSize:(i+1)*HiddenSize]
	FTBias    []float32 // [HiddenSize]

	L1Weights []float32 // [2*HiddenSize][L2Size]
	L1Bias    []float32 // [L2Size]

	L2Weights []float32 // [L2Size][L3Size]
	L2Bias    []float32 // [L3Size]

	L3Weights []float32 // [L3Size]
	L3Bias    float32
}

// currently loaded network
var Net *Network

func addFeature(acc *[HiddenSize]float32, feature int) {
	row := Net.FTWeights[feature*HiddenSize : feature*HiddenSize+HiddenSize]
	for i := range HiddenSize {
		acc[i] += row[i]
	}
}

func removeFeature(acc *[HiddenSize]float32, feature int) {
	row := Net.FTWeights[feature*HiddenSize : feature*HiddenSize+HiddenSize]
	for i := range HiddenSize {
		acc[i] -= row[i]
	}
}

func clippedReLU(x float32) float32 {
	if x < 0 {
		return 0
	}
	return x
}

// runs forward pass and returns centipawn score from perspective of side to move
func Evaluate(b *board.Board) int {
	stm := int(b.Side)
	other := 1 - stm

	var input [2 * HiddenSize]float32
	for i := range HiddenSize {
		input[i] = clippedReLU(b.NNUEAcc[stm][i])
		input[HiddenSize+i] = clippedReLU(b.NNUEAcc[other][i])
	}

	var l2 [L2Size]float32
	for j := range L2Size {
		sum := Net.L1Bias[j]
		for i := range 2 * HiddenSize {
			sum += input[i] * Net.L1Weights[i*L2Size+j]
		}
		l2[j] = clippedReLU(sum)
	}

	var l3 [L3Size]float32
	for j := range L3Size {
		sum := Net.L2Bias[j]
		for i := range L2Size {
			sum += l2[i] * Net.L2Weights[i*L3Size+j]
		}
		l3[j] = clippedReLU(sum)
	}

	out := Net.L3Bias
	for i := range L3Size {
		out += l3[i] * Net.L3Weights[i]
	}

	return int(out)
}
