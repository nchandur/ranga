package main

type CastleBit int

const (
	WKSide CastleBit = 1 << iota
	WQSide
	BKSide
	BQSide
)
