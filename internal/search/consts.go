package search

const (
	MAX_DEPTH        int = 8       // defines maximum search depth in plies
	MAX_PLY          int = 64      // defines the maximum distance from the root node
	MAX_HISTORY      int = 1 << 14 // cap for history bonus
	ISMATE           int = 1000000 // base score used to identify mate condition
	MATESCORE        int = 900000  // threshold above which scores represent guaranteed checkmate sequences within x plies
	INFINITY         int = 1200000 // upper numeric boundary
	FULL_DEPTH_MOVES int = 4       // number of moves evaluated at full depth before LMR
	REDUCTION_LIMIT  int = 3       // minimum remaining search depth required to trigger LMR
)

// assigns higher priorities to captures where lower-value piece captures higher-value piece,
var MVVLVA = [12][12]int{
	{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
	{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
	{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
	{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
	{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
	{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},

	{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
	{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
	{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
	{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
	{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
	{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
}
