package main

import "math/rand"

var PieceKeys [13][120]uint64
var SideKey uint64
var CastleKeys [16]uint64
var EnPassKeys [120]uint64

func initHashkeys() {
	rng := rand.New(rand.NewSource(19991211))

	for piece := range 13 {
		for sq := range 120 {
			PieceKeys[piece][sq] = rng.Uint64()
		}
	}

	for i := range 16 {
		CastleKeys[i] = rng.Uint64()
	}

	for sq := range 120 {
		EnPassKeys[sq] = rng.Uint64()
	}

	SideKey = rng.Uint64()
}
