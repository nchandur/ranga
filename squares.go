package main

import "fmt"

type Square int

const (

	// First Rank
	A1 Square = 21
	B1 Square = 22
	C1 Square = 23
	D1 Square = 24
	E1 Square = 25
	F1 Square = 26
	G1 Square = 27
	H1 Square = 28

	// Last Rank
	A8 Square = 91
	B8 Square = 92
	C8 Square = 93
	D8 Square = 94
	E8 Square = 95
	F8 Square = 96
	G8 Square = 97
	H8 Square = 98

	// Special Squares
	NoSquare Square = 99
	Offboard Square = 120
)

func (s *Square) String() string {
	return fmt.Sprintf("%c%d", 'a'+FilesBoard[*s], RanksBoard[*s]+1)
}
