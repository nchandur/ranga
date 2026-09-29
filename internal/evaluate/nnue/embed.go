package nnue

import _ "embed"

//go:embed weights/net.bin
var embeddedWeights []byte

// parses the network compiled directly into the binary
func LoadEmbedded() (*Network, error) {
	return parseNetwork(embeddedWeights)
}
