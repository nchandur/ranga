package main

import "math/rand"

var PieceKeys [130]uint64
var SideKey uint64
var CastleKeys [16]uint64

func initHashkeys() {
	rng := rand.New(rand.NewSource(7867))

	for sq := range 130 {
		PieceKeys[sq] = rng.Uint64()
	}

	for i := range 16 {
		CastleKeys[i] = rng.Uint64()
	}

	SideKey = rng.Uint64()
}
