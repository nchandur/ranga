package main

// returns square index for a given file and rank
func FR2Square(file File, rank Rank) Square {
	return Square((21 + int(file)) + (int(rank) * 10))
}

// compare slices out of order
func CompareUnordered(a, b []int) bool {

	if len(a) != len(b) {
		return false
	}

	count := make(map[int]int)

	for _, v := range a {
		count[v]++
	}

	for _, v := range b {
		count[v]--

		if count[v] < 0 {
			return false
		}
	}

	return true

}
