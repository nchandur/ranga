package search

const (
	MAX_DEPTH        int = 8       // defines maximum search depth in plies
	MAX_PLY          int = 256     // defines the maximum distance from the root node
	MAX_HISTORY      int = 1 << 14 // cap for history bonus
	ISMATE           int = 1000000 // base score used to identify mate condition
	MATESCORE        int = 900000  // threshold above which scores represent guaranteed checkmate sequences within x plies
	INFINITY         int = 1200000 // upper numeric boundary
	FULL_DEPTH_MOVES int = 4       // number of moves evaluated at full depth before LMR
	REDUCTION_LIMIT  int = 3       // minimum remaining search depth required to trigger LMR
	BIG_DELTA        int = 200
	DEFAULT_TT_SIZE  int = 16
)

const (
	aspirationMinDepth = 5
	aspirationDelta    = 20
	aspirationMaxDelta = 800 // beyond this, window fully open on that side
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

// futility margin
var FutilityMargin [4]int = [4]int{0, 100, 300, 500}

// the most threads the Threads option allows
const MAX_THREADS = 8

// how often a searcher publishes its node count and looks at the stop flag
const pollInterval = 2048

// helpers spread over different depths instead of repeating the main thread's work
var (
	skipSize  = [...]int{1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	skipPhase = [...]int{0, 1, 0, 1, 2, 3, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7}
)

const helperMaxDepth = 100
