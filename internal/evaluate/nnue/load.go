package nnue

import (
	"encoding/binary"
	"fmt"
	"os"
)

// feature weights + feature bias + (output weights + output bias) per bucket
const expectedNetworkBytes = (768*HiddenSize + HiddenSize + OutputBuckets*2*HiddenSize + OutputBuckets) * 2

func LoadNetwork(path string) (*Network, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading network file: %w", err)
	}
	return parseNetwork(data)
}

func parseNetwork(data []byte) (*Network, error) {
	if len(data) < expectedNetworkBytes {
		return nil, fmt.Errorf("network data too small: got %d bytes, want at least %d", len(data), expectedNetworkBytes)
	}

	r := bytesReader(data[:expectedNetworkBytes])
	var net Network

	for f := range net.FeatureWeights {
		for i := range HiddenSize {
			net.FeatureWeights[f].Values[i] = int16(binary.LittleEndian.Uint16(r.next(2)))
		}
	}
	for i := range HiddenSize {
		net.FeatureBias.Values[i] = int16(binary.LittleEndian.Uint16(r.next(2)))
	}

	// output weights are bucket-major]
	for b := range OutputBuckets {
		for i := range net.OutputWeights[b] {
			net.OutputWeights[b][i] = int16(binary.LittleEndian.Uint16(r.next(2)))
		}
	}
	// followed by one bias per bucket
	for b := range OutputBuckets {
		net.OutputBias[b] = int16(binary.LittleEndian.Uint16(r.next(2)))
	}

	return &net, nil
}

type byteCursor struct {
	data []byte
	pos  int
}

func bytesReader(data []byte) *byteCursor { return &byteCursor{data: data} }

func (c *byteCursor) next(n int) []byte {
	b := c.data[c.pos : c.pos+n]
	c.pos += n
	return b
}
