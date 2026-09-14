package nnue

import (
	"encoding/binary"
	"fmt"
	"os"
)

const expectedNetworkBytes = (768*HiddenSize + HiddenSize + 2*HiddenSize + 1) * 2

func LoadNetwork(path string) (*Network, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading network file: %w", err)
	}
	if len(data) < expectedNetworkBytes {
		return nil, fmt.Errorf("network file too small: got %d bytes, want at least %d", len(data), expectedNetworkBytes)
	}

	r := bytesReader(data[:expectedNetworkBytes])
	var net Network

	for f := range net.FeatureWeights {
		for i := range HiddenSize {
			net.FeatureWeights[f].Values[i] = int(int16(binary.LittleEndian.Uint16(r.next(2))))
		}
	}
	for i := range HiddenSize {
		net.FeatureBias.Values[i] = int(int16(binary.LittleEndian.Uint16(r.next(2))))
	}
	for i := range net.OutputWeights {
		net.OutputWeights[i] = int(int16(binary.LittleEndian.Uint16(r.next(2))))
	}
	net.OutputBias = int(int16(binary.LittleEndian.Uint16(r.next(2))))

	return &net, nil
}

// tiny helper to walk the byte slice without a full bufio.Reader
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
