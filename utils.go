package main

// returns square index for a given file and rank
func FR2Square(file File, rank Rank) Square {
	return Square((21 + int(file)) + (int(rank) * 10))
}
