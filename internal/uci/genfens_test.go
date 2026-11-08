package uci

import (
	"fmt"
	"math/bits"
	"ranga/internal/board"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Test fixtures
// ---------------------------------------------------------------------------

type startCase struct {
	name     string
	fen      string
	standard bool // reached from the standard start; enables "no promotion in <=10 plies" checks
}

var startCases = []startCase{
	{"startpos", board.START, true},
	{"kiwipete", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", false},
	{"black-to-move-ep", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", false},
	{"ep-available", "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3", false},
	{"castling-only", "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", false},
	{"rook-endgame", "8/5pk1/6p1/8/8/6P1/5PK1/R7 w - - 0 40", false},
	{"promotion-race", "8/1P4k1/8/8/8/8/4K1p1/8 w - - 0 1", false},
}

var pliesList = []int{0, 1, 2, 3, 7, 8, 9}

// seeds exercise all 64 bits (OpenBench puts the workload ID in the upper half)
var testSeeds = []uint64{
	0, 1, 2, 3, 42, 12345,
	1 << 32, 1<<32 | 1, 0xFFFFFFFF,
	0xFFFFFFFF00000000, 0xFFFFFFFFFFFFFFFF, 0xDEADBEEFCAFEBABE,
}

const testCount = 25

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mustParse(t *testing.T, fen string) board.Board {
	t.Helper()
	b := board.NewBoard()
	if err := b.ParseFEN(fen); err != nil {
		t.Fatalf("bad test FEN %q: %v", fen, err)
	}
	return b
}

func collect(t *testing.T, start *board.Board, count, plies int, seed uint64) []board.Board {
	t.Helper()
	var out []board.Board
	n := genPositions(start, count, plies, seed, func(b *board.Board) {
		out = append(out, b.Preserve())
	})
	if n != len(out) {
		t.Fatalf("genPositions returned %d but emitted %d positions", n, len(out))
	}
	return out
}

func fensOf(bs []board.Board) []string {
	res := make([]string, len(bs))
	for i := range bs {
		res[i] = bs[i].FEN()
	}
	return res
}

func isWhitePiece(p board.Piece) bool { return p >= board.WP && p <= board.WK }
func isBlackPiece(p board.Piece) bool { return p >= board.BP && p <= board.BK }

func bitSet(bb board.BitBoard, sq int) bool { return uint64(bb)>>uint(sq)&1 == 1 }

const (
	iP = iota
	iN
	iB
	iR
	iQ
	iK
)

func pieceSet(white bool) [6]board.Piece {
	if white {
		return [6]board.Piece{board.WP, board.WN, board.WB, board.WR, board.WQ, board.WK}
	}
	return [6]board.Piece{board.BP, board.BN, board.BB, board.BR, board.BQ, board.BK}
}

// counts returns pawn, knight, bishop, rook, queen, king counts, from the mailbox
func counts(b *board.Board, white bool) [6]int {
	set := pieceSet(white)
	var c [6]int
	for _, pc := range b.Mailbox {
		for i, want := range set {
			if pc == want {
				c[i]++
			}
		}
	}
	return c
}

func nonKing(c [6]int) int { return c[iP] + c[iN] + c[iB] + c[iR] + c[iQ] }

func kingSquare(b *board.Board, white bool) (int, bool) {
	k := board.BK
	if white {
		k = board.WK
	}
	for sq := 0; sq < 64; sq++ {
		if b.Mailbox[sq] == k {
			return sq, true
		}
	}
	return 0, false
}

// attackedBy is an independent, mailbox-only attack detector (row-major, a8 = 0).
// It deliberately shares no code with the engine's bitboard attack tables.
func attackedBy(b *board.Board, r, f int, byWhite bool) bool {
	pick := func(w, bl board.Piece) board.Piece {
		if byWhite {
			return w
		}
		return bl
	}
	P, N, B := pick(board.WP, board.BP), pick(board.WN, board.BN), pick(board.WB, board.BB)
	R, Q, K := pick(board.WR, board.BR), pick(board.WQ, board.BQ), pick(board.WK, board.BK)

	at := func(r, f int) (board.Piece, bool) {
		if r < 0 || r > 7 || f < 0 || f > 7 {
			return board.Empty, false
		}
		return b.Mailbox[r*8+f], true
	}

	// a white pawn on (r+1, f±1) attacks (r, f); a black pawn on (r-1, f±1) does
	pr := r + 1
	if !byWhite {
		pr = r - 1
	}
	for _, df := range []int{-1, 1} {
		if p, ok := at(pr, f+df); ok && p == P {
			return true
		}
	}

	for _, d := range [][2]int{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}} {
		if p, ok := at(r+d[0], f+d[1]); ok && p == N {
			return true
		}
	}

	for dr := -1; dr <= 1; dr++ {
		for df := -1; df <= 1; df++ {
			if dr == 0 && df == 0 {
				continue
			}
			if p, ok := at(r+dr, f+df); ok && p == K {
				return true
			}
		}
	}

	slide := func(dirs [][2]int, a, c board.Piece) bool {
		for _, d := range dirs {
			rr, ff := r+d[0], f+d[1]
			for {
				p, ok := at(rr, ff)
				if !ok {
					break
				}
				if p != board.Empty {
					if p == a || p == c {
						return true
					}
					break
				}
				rr += d[0]
				ff += d[1]
			}
		}
		return false
	}
	if slide([][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}, R, Q) {
		return true
	}
	return slide([][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}, B, Q)
}

// walk visits every leaf of the legal-move tree to the given depth
func walk(b board.Board, depth int, visit func(*board.Board)) {
	if depth == 0 {
		visit(&b)
		return
	}
	list := board.NewMoveList()
	list.GenerateMoves(&b)
	for i := range list.Count {
		c := b.Preserve()
		if c.MakeMove(list.Moves[i], false) {
			walk(c, depth-1, visit)
		}
	}
}

// key4 is the position identity ignoring move counters
func key4(b *board.Board) string {
	return strings.Join(strings.Fields(b.FEN())[:4], " ")
}

// ---------------------------------------------------------------------------
// Individual validity checks (each returns nil when the position is fine)
// ---------------------------------------------------------------------------

// internal state agrees with itself: mailbox, bitboards, occupancies
// (the Zobrist key is checked separately in TestZobristKeyMatchesRecompute)
func checkConsistency(b *board.Board) error {
	var white, black uint64
	for sq := 0; sq < 64; sq++ {
		pc := b.Mailbox[sq]
		for q := board.WP; q <= board.BK; q++ {
			if has := bitSet(b.PieceBitBoards[q], sq); has != (pc == q) {
				return fmt.Errorf("square %d: mailbox=%v but bitboard[%v] bit=%v", sq, pc, q, has)
			}
		}
		if isWhitePiece(pc) {
			white |= 1 << uint(sq)
		} else if isBlackPiece(pc) {
			black |= 1 << uint(sq)
		}
	}
	if uint64(b.Occupancies[board.White]) != white {
		return fmt.Errorf("white occupancy mismatch: have %#x want %#x", uint64(b.Occupancies[board.White]), white)
	}
	if uint64(b.Occupancies[board.Black]) != black {
		return fmt.Errorf("black occupancy mismatch: have %#x want %#x", uint64(b.Occupancies[board.Black]), black)
	}
	if uint64(b.Occupancies[board.Both]) != white|black {
		return fmt.Errorf("combined occupancy mismatch: have %#x want %#x", uint64(b.Occupancies[board.Both]), white|black)
	}
	if bits.OnesCount64(white&black) != 0 {
		return fmt.Errorf("white and black occupancy overlap")
	}
	if b.Side != board.White && b.Side != board.Black {
		return fmt.Errorf("invalid side to move %v", b.Side)
	}
	return nil
}

// the FEN is well formed, ParseFEN accepts it, and it reproduces the same position
func checkFENRoundTrip(b *board.Board) error {
	fen := b.FEN()
	if n := len(strings.Fields(fen)); n != 6 {
		return fmt.Errorf("FEN has %d fields, want 6", n)
	}
	p := board.NewBoard()
	if err := p.ParseFEN(fen); err != nil {
		return fmt.Errorf("ParseFEN rejected own FEN: %v", err)
	}
	if got := p.FEN(); got != fen {
		return fmt.Errorf("round trip changed FEN to %q", got)
	}
	switch {
	case p.Mailbox != b.Mailbox:
		return fmt.Errorf("round trip changed piece placement")
	case p.Side != b.Side:
		return fmt.Errorf("round trip changed side to move")
	case p.Castle != b.Castle:
		return fmt.Errorf("round trip changed castling rights")
	case p.EnPassant != b.EnPassant:
		return fmt.Errorf("round trip changed en passant square")
	case p.FiftyMove != b.FiftyMove:
		return fmt.Errorf("round trip changed halfmove clock")
	case p.Ply != b.Ply:
		return fmt.Errorf("round trip changed ply: %d vs %d", p.Ply, b.Ply)
	}
	return nil
}

// piece-count sanity: one king each, no pawns on back ranks, promotions fit in 8 pawns
func checkStructure(b *board.Board) error {
	for sq := 0; sq < 8; sq++ {
		for _, idx := range []int{sq, 56 + sq} {
			if pc := b.Mailbox[idx]; pc == board.WP || pc == board.BP {
				return fmt.Errorf("pawn on back rank (square %d)", idx)
			}
		}
	}
	for _, white := range []bool{true, false} {
		name := map[bool]string{true: "white", false: "black"}[white]
		c := counts(b, white)
		if c[iK] != 1 {
			return fmt.Errorf("%s has %d kings", name, c[iK])
		}
		if c[iP] > 8 {
			return fmt.Errorf("%s has %d pawns", name, c[iP])
		}
		if c[iP]+c[iN]+c[iB]+c[iR]+c[iQ]+c[iK] > 16 {
			return fmt.Errorf("%s has more than 16 pieces", name)
		}
		excess := max(0, c[iQ]-1) + max(0, c[iR]-2) + max(0, c[iB]-2) + max(0, c[iN]-2)
		if c[iP]+excess > 8 {
			return fmt.Errorf("%s has %d pawns plus %d promoted pieces (>8)", name, c[iP], excess)
		}
	}
	return nil
}

// the side NOT to move must not be in check (such a position is unreachable)
func checkOpponentNotInCheck(b *board.Board) error {
	stmWhite := b.Side == board.White
	sq, ok := kingSquare(b, !stmWhite)
	if !ok {
		return fmt.Errorf("opponent king missing")
	}
	if attackedBy(b, sq/8, sq%8, stmWhite) {
		return fmt.Errorf("side not to move is in check (independent detector)")
	}
	color := board.Black
	if stmWhite {
		color = board.White
	}
	if b.IsSquareAttacked(board.Square(sq), color) {
		return fmt.Errorf("side not to move is in check (engine detector)")
	}
	return nil
}

// the engine's bitboard attack code agrees with the independent detector on all squares
func checkAttackAgreement(b *board.Board) error {
	for sq := 0; sq < 64; sq++ {
		for _, byWhite := range []bool{true, false} {
			color := board.Black
			if byWhite {
				color = board.White
			}
			want := attackedBy(b, sq/8, sq%8, byWhite)
			if got := b.IsSquareAttacked(board.Square(sq), color); got != want {
				return fmt.Errorf("square %d attacked-by-white=%v: engine=%v independent=%v", sq, byWhite, got, want)
			}
		}
	}
	return nil
}

// the side to move must have at least one legal move (not mate or stalemate)
func checkNotTerminal(b *board.Board) error {
	list := board.NewMoveList()
	list.GenerateMoves(b)
	for i := range list.Count {
		c := b.Preserve()
		if c.MakeMove(list.Moves[i], false) {
			return nil
		}
	}
	return fmt.Errorf("side to move has no legal moves (mate or stalemate)")
}

// castling rights require the king and rook to be on their home squares
func checkCastlingRights(b *board.Board) error {
	need := func(has bool, name string, ksq, rsq board.Square, kp, rp board.Piece) error {
		if has && (b.Mailbox[ksq] != kp || b.Mailbox[rsq] != rp) {
			return fmt.Errorf("castling right %s set but king/rook not on home squares", name)
		}
		return nil
	}
	for _, err := range []error{
		need(b.Castle&board.WKCA != 0, "K", board.E1, board.H1, board.WK, board.WR),
		need(b.Castle&board.WQCA != 0, "Q", board.E1, board.A1, board.WK, board.WR),
		need(b.Castle&board.BKCA != 0, "k", board.E8, board.H8, board.BK, board.BR),
		need(b.Castle&board.BQCA != 0, "q", board.E8, board.A8, board.BK, board.BR),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}

// en passant square must be consistent with a pawn having just double-pushed
func checkEnPassant(b *board.Board) error {
	if b.EnPassant == board.NoSquare {
		return nil
	}
	ep := int(b.EnPassant)
	if ep < 0 || ep > 63 {
		return fmt.Errorf("en passant square %d out of range", ep)
	}
	if b.Mailbox[ep] != board.Empty {
		return fmt.Errorf("en passant square %d is occupied", ep)
	}
	if b.Side == board.White {
		// black just pushed: ep on rank 6 (row 2), pawn on rank 5 (row 3), origin rank 7 (row 1) empty
		if ep/8 != 2 || b.Mailbox[ep+8] != board.BP || b.Mailbox[ep-8] != board.Empty {
			return fmt.Errorf("inconsistent en passant square %d for white to move", ep)
		}
	} else {
		// white just pushed: ep on rank 3 (row 5), pawn on rank 4 (row 4), origin rank 2 (row 6) empty
		if ep/8 != 5 || b.Mailbox[ep-8] != board.WP || b.Mailbox[ep+8] != board.Empty {
			return fmt.Errorf("inconsistent en passant square %d for black to move", ep)
		}
	}
	return nil
}

// counters and side to move must reflect exactly `plies` plies from the start
func counterCheck(start *board.Board, plies int) func(*board.Board) error {
	return func(b *board.Board) error {
		if b.Ply != start.Ply+plies {
			return fmt.Errorf("ply = %d, want %d", b.Ply, start.Ply+plies)
		}
		wantSide := start.Side
		if plies%2 == 1 {
			wantSide ^= 1
		}
		if b.Side != wantSide {
			return fmt.Errorf("side to move = %v, want %v", b.Side, wantSide)
		}
		if b.FiftyMove < 0 || b.FiftyMove > start.FiftyMove+plies {
			return fmt.Errorf("halfmove clock %d impossible after %d plies from %d", b.FiftyMove, plies, start.FiftyMove)
		}
		if b.Castle&^start.Castle != 0 {
			return fmt.Errorf("castling rights appeared that the start position did not have")
		}
		return nil
	}
}

// necessary conditions for reachability within `plies` plies (material can only be lost,
// each capture needs a move by the opponent, no promotions when too few moves are available)
func reachabilityBounds(start *board.Board, plies int, noPromotion bool) func(*board.Board) error {
	stmMoves, otherMoves := (plies+1)/2, plies/2
	whiteMoves, blackMoves := stmMoves, otherMoves
	if start.Side == board.Black {
		whiteMoves, blackMoves = otherMoves, stmMoves
	}
	return func(b *board.Board) error {
		for _, white := range []bool{true, false} {
			before, after := counts(start, white), counts(b, white)
			enemyMoves := blackMoves
			if !white {
				enemyMoves = whiteMoves
			}
			lost := nonKing(before) - nonKing(after)
			if lost < 0 || lost > enemyMoves {
				return fmt.Errorf("white=%v lost %d pieces, but the opponent only has %d moves", white, lost, enemyMoves)
			}
			if after[iP] > before[iP] {
				return fmt.Errorf("white=%v gained pawns (%d -> %d)", white, before[iP], after[iP])
			}
			if noPromotion {
				for i := iN; i <= iQ; i++ {
					if after[i] > before[i] {
						return fmt.Errorf("white=%v piece type %d count rose %d -> %d without time to promote", white, i, before[i], after[i])
					}
				}
			}
		}
		return nil
	}
}

type namedCheck struct {
	name string
	fn   func(*board.Board) error
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// The main property test: every generated position passes every validity check.
func TestGeneratedPositionsAreValid(t *testing.T) {
	for _, tc := range startCases {
		start := mustParse(t, tc.fen)

		for _, plies := range pliesList {
			t.Run(fmt.Sprintf("%s/plies=%d", tc.name, plies), func(t *testing.T) {
				checks := []namedCheck{
					{"internal consistency", checkConsistency},
					{"FEN round trip", checkFENRoundTrip},
					{"piece structure", checkStructure},
					{"opponent not in check", checkOpponentNotInCheck},
					{"attack detector agreement", checkAttackAgreement},
					{"side to move not mated/stalemated", checkNotTerminal},
					{"castling rights", checkCastlingRights},
					{"en passant square", checkEnPassant},
					{"counters", counterCheck(&start, plies)},
					{"reachability bounds", reachabilityBounds(&start, plies, tc.standard && plies <= 10)},
				}

				for _, seed := range testSeeds {
					positions := collect(t, &start, testCount, plies, seed)
					if len(positions) != testCount {
						t.Fatalf("seed %#x: got %d positions, want %d", seed, len(positions), testCount)
					}
					for i := range positions {
						for _, c := range checks {
							if err := c.fn(&positions[i]); err != nil {
								t.Fatalf("%s failed\n  start: %s\n  plies: %d  seed: %#x\n  fen:   %s\n  error: %v",
									c.name, tc.fen, plies, seed, positions[i].FEN(), err)
							}
						}
					}
				}
			})
		}
	}
}

// Grounds the reachability oracle below: the engine's move generator must match known perft values.
func TestPerftSanity(t *testing.T) {
	cases := []struct {
		name  string
		fen   string
		nodes []int // perft(1), perft(2), ...
	}{
		{"startpos", board.START, []int{20, 400, 8902, 197281}},
		{"kiwipete", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", []int{48, 2039, 97862}},
	}
	for _, tc := range cases {
		start := mustParse(t, tc.fen)
		for d, want := range tc.nodes {
			got := 0
			walk(start, d+1, func(*board.Board) { got++ })
			if got != want {
				t.Errorf("%s perft(%d) = %d, want %d: movegen/MakeMove is broken, so genfens results cannot be trusted",
					tc.name, d+1, got, want)
			}
		}
	}
}

// Exact reachability for small depths: every generated position must be a leaf of the
// legal-move tree from the start position at exactly that depth.
func TestGeneratedPositionsAreReachable(t *testing.T) {
	for _, tc := range startCases {
		start := mustParse(t, tc.fen)

		maxDepth := 3
		if tc.standard {
			maxDepth = 4
		}
		for plies := 1; plies <= maxDepth; plies++ {
			t.Run(fmt.Sprintf("%s/plies=%d", tc.name, plies), func(t *testing.T) {
				reachable := map[string]struct{}{}
				walk(start, plies, func(b *board.Board) { reachable[key4(b)] = struct{}{} })

				for _, seed := range testSeeds {
					for _, b := range collect(t, &start, testCount, plies, seed) {
						if _, ok := reachable[key4(&b)]; !ok {
							t.Fatalf("unreachable position\n  start: %s\n  plies: %d  seed: %#x\n  fen:   %s",
								tc.fen, plies, seed, b.FEN())
						}
					}
				}
			})
		}
	}
}

// Same seed and start position must give identical output; the generator must use all 64 seed bits.
func TestGenPositionsDeterministicAndSeedSensitive(t *testing.T) {
	start := mustParse(t, board.START)
	const seed = uint64(0xDEADBEEF00000042)

	a := fensOf(collect(t, &start, 50, genfensPlies, seed))
	b := fensOf(collect(t, &start, 50, genfensPlies, seed))
	if strings.Join(a, "\n") != strings.Join(b, "\n") {
		t.Fatal("same seed produced different output")
	}

	for name, other := range map[string]uint64{
		"lower 32 bits differ": seed ^ 1,
		"upper 32 bits differ": seed ^ (1 << 32),
	} {
		c := fensOf(collect(t, &start, 50, genfensPlies, other))
		if strings.Join(a, "\n") == strings.Join(c, "\n") {
			t.Errorf("%s: output identical to the original seed", name)
		}
	}
}

// Generation must not alter the position it starts from.
func TestStartPositionNotMutated(t *testing.T) {
	for _, tc := range startCases {
		start := mustParse(t, tc.fen)
		snapshot := start.Preserve()
		for _, seed := range testSeeds {
			collect(t, &start, 10, genfensPlies, seed)
		}
		if start != snapshot {
			t.Errorf("%s: start position was modified by genPositions", tc.name)
		}
	}
}

// Exactly the requested number of positions is produced.
func TestGeneratedCount(t *testing.T) {
	start := mustParse(t, board.START)
	for _, count := range []int{0, 1, 7, 100} {
		if got := len(collect(t, &start, count, genfensPlies, 5)); got != count {
			t.Errorf("count=%d: got %d positions", count, got)
		}
	}
}

// Positions should actually be random: mostly distinct, and not just the start position.
func TestGeneratedDiversity(t *testing.T) {
	start := mustParse(t, board.START)
	startFEN := start.FEN()

	total, atStart := 0, 0
	distinct := map[string]struct{}{}
	for _, seed := range testSeeds {
		for _, fen := range fensOf(collect(t, &start, 100, genfensPlies, seed)) {
			total++
			distinct[fen] = struct{}{}
			if fen == startFEN {
				atStart++
			}
		}
	}
	ratio := float64(len(distinct)) / float64(total)
	t.Logf("%d positions, %d distinct (%.1f%%), %d equal to start", total, len(distinct), 100*ratio, atStart)

	if ratio < 0.80 {
		t.Errorf("only %.1f%% of generated positions are distinct", 100*ratio)
	}
	if float64(atStart) > 0.01*float64(total) {
		t.Errorf("%d of %d positions equal the start position", atStart, total)
	}
}

// If the start position is already mate/stalemate, generation must terminate and emit nothing.
func TestTerminalStartTerminates(t *testing.T) {
	cases := map[string]string{
		"stalemate": "7k/5Q2/6K1/8/8/8/8/8 b - - 0 1",
		"checkmate": "R5k1/5ppp/8/8/8/8/8/6K1 b - - 0 1",
	}
	for name, fen := range cases {
		t.Run(name, func(t *testing.T) {
			start := mustParse(t, fen)
			done := make(chan int, 1)
			go func() {
				done <- genPositions(&start, 5, genfensPlies, 1, func(*board.Board) {})
			}()
			select {
			case n := <-done:
				if n != 0 {
					t.Fatalf("emitted %d positions from a terminal start position", n)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("genPositions did not terminate on a terminal start position")
			}
		})
	}
}

// The incrementally updated Zobrist key must equal a from-scratch recomputation.
// This is an engine invariant (transposition table, repetition detection) that is
// independent of the FENs genfens prints, so it is kept separate from the validity checks.
func TestZobristKeyMatchesRecompute(t *testing.T) {
	for _, tc := range startCases {
		start := mustParse(t, tc.fen)
		for _, plies := range pliesList {
			t.Run(fmt.Sprintf("%s/plies=%d", tc.name, plies), func(t *testing.T) {
				for _, seed := range testSeeds {
					for _, b := range collect(t, &start, testCount, plies, seed) {
						want := b.GenerateKey()
						if b.Key == want {
							continue
						}
						hint := ""
						if b.Key^want == board.SideKey {
							hint = "\n  the keys differ by exactly SideKey: GenerateKey and MakeMove disagree on the side-to-move term"
						}
						t.Fatalf("incremental key %#x != recomputed key %#x\n  start: %s\n  plies: %d  seed: %#x\n  fen:   %s%s",
							b.Key, want, tc.fen, plies, seed, b.FEN(), hint)
					}
				}
			})
		}
	}
}
