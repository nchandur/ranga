package nnue

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
)

// builds a network with small random weights
func RandomNetwork() *Network {
	return &Network{
		FTWeights: randSlice(InputSize * HiddenSize),
		FTBias:    randSlice(HiddenSize),
		L1Weights: randSlice(2 * HiddenSize * L2Size),
		L1Bias:    randSlice(L2Size),
		L2Weights: randSlice(L2Size * L3Size),
		L2Bias:    randSlice(L3Size),
		L3Weights: randSlice(L3Size),
		L3Bias:    (rand.Float32()*2 - 1) * 0.01,
	}
}

func randSlice(n int) []float32 {
	s := make([]float32, n)
	for i := range s {
		s[i] = (rand.Float32()*2 - 1) * 0.01
	}
	return s
}

// reads weights from a simple custom binary format
func LoadNetwork(path string) (*Network, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var magic [4]byte
	if _, err := f.Read(magic[:]); err != nil {
		return nil, err
	}
	if string(magic[:]) != "RNUE" {
		return nil, fmt.Errorf("nnue: bad magic %q, expected \"RNUE\"", magic)
	}

	n := &Network{}
	sections := []struct {
		dst *[]float32
		len int
	}{
		{&n.FTWeights, InputSize * HiddenSize},
		{&n.FTBias, HiddenSize},
		{&n.L1Weights, 2 * HiddenSize * L2Size},
		{&n.L1Bias, L2Size},
		{&n.L2Weights, L2Size * L3Size},
		{&n.L2Bias, L3Size},
		{&n.L3Weights, L3Size},
	}
	for _, s := range sections {
		buf := make([]float32, s.len)
		if err := binary.Read(f, binary.LittleEndian, buf); err != nil {
			return nil, fmt.Errorf("nnue: reading weights: %w", err)
		}
		*s.dst = buf
	}
	if err := binary.Read(f, binary.LittleEndian, &n.L3Bias); err != nil {
		return nil, fmt.Errorf("nnue: reading final bias: %w", err)
	}
	return n, nil
}
