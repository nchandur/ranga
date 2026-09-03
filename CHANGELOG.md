# Changelog

All notable changes to this project will be documented in this file.

## [v1.11] - 2026-09-03

### Added
* feat(uci): support `go nodes <n>` to limit search by node count [[#36](https://github.com/nchandur/ranga/pull/36)]
* feat(time-control): soft/hard limits for move allocation [[#37](https://github.com/nchandur/ranga/pull/37)]

## [v1.10] - 2026-09-02

### Added
* feat(evaluate): tapered evaluation [[#34](https://github.com/nchandur/ranga/pull/34)]

## [v1.9] - 2026-09-02 

### Added
* feat(evaluate): piece-square tables, mobility bonuses/penalties and king safety [[#32](https://github.com/nchandur/ranga/pull/32)]
* chore(build): Go version bump to `1.27.1` [[#32](https://github.com/nchandur/ranga/pull/32)]

## [v1.8] - 2026-09-02

### Added
* feat(search): null move pruning [[#30](https://github.com/nchandur/ranga/pull/30)]

### Fixed
* fix(search): no LMR on first move [[#29](https://github.com/nchandur/ranga/pull/29)]
* fix(search): correct transposition table replacement scheme [[#29](https://github.com/nchandur/ranga/pull/29)]
* fix(search): store checkmate/stalemate scores in transposition table [[#29](https://github.com/nchandur/ranga/pull/29)]

## [v1.7] - 2026-09-01
* feat(search): late move reduction (LMR) with PVS re-search [[#26](https://github.com/nchandur/ranga/pull/26)]

## [v1.6] - 2026-08-27 

### Added
* feat(search): killer moves and history heuristics [[#20](https://github.com/nchandur/ranga/pull/20)]

### Fixed
* fix(search): pv root node corruption [[#21](https://github.com/nchandur/ranga/pull/21)]
* fix(search): history heuristic on strictly legal quiet moves [[#21](https://github.com/nchandur/ranga/pull/21)]
* fix(time-management): time crash vulnerability [[#21](https://github.com/nchandur/ranga/pull/21)]

## [v1.5] - 2026-08-25

### Added
* feat(search): transposition table for storing previously searched positions [[#17](https://github.com/nchandur/ranga/pull/17)]

### Fixed
* fix(search): enforce minimum 1MB transposition table size [[#18](https://github.com/nchandur/ranga/pull/18)]
* fix(search): root node now probes TT for move ordering instead of using no move [[#18](https://github.com/nchandur/ranga/pull/18)]
* fix(search): TT replacement scheme now protects same-key entries from shallower overwrites [[#18](https://github.com/nchandur/ranga/pull/18)]

## [v1.4] - 2026-08-25

### Added
* feat(search): repetition detection and 50-move rule [[#15](https://github.com/nchandur/ranga/pull/15)]

### Fixed
* fix(search): correct ply count tracking during search traversal [[#15](https://github.com/nchandur/ranga/pull/15)]

## [v1.3] - 2026-08-23

### Added
* feat(search): principal variation table [[#12](https://github.com/nchandur/ranga/pull/12)]

### Fixed
* fix(board): correct side-to-move in log [[#12](https://github.com/nchandur/ranga/pull/12)]

## [v1.2] - 2026-08-23

### Added
* feat(move-order): MVVLVA [[#10](https://github.com/nchandur/ranga/pull/10)]

## [v1.1] - 2026-08-23

### Added
* feat(search): quiescence search [[#8](https://github.com/nchandur/ranga/pull/8)]

### Fixed
* fix(uci): correctly parse version number in UCI handshake

## [v1.0] - 2026-08-20

### Added

* feat(time-control): basic time control [[#4](https://github.com/nchandur/ranga/pull/4)]
* feat(uci): basic uci commands [[#4](https://github.com/nchandur/ranga/pull/4)]
* feat(uci): search info reporting (depth, score, nodes) [[#4](https://github.com/nchandur/ranga/pull/4)]
* feat(search): iterative deepening [[#4](https://github.com/nchandur/ranga/pull/4)]
* feat(evaluate): static piece value evaluation [[#3](https://github.com/nchandur/ranga/pull/3)]
* feat(search): basic alpha-beta search [[#3](https://github.com/nchandur/ranga/pull/3)]
* feat(board): move generation [[#2](https://github.com/nchandur/ranga/pull/2)]
* feat(board): FEN parsing, utilities, and Zobrist hashing for position signatures [[#1](https://github.com/nchandur/ranga/pull/1)]
* feat(board): pre-computed attack tables, including 1D arrays for leapers and magic bitboards for sliders [[#1](https://github.com/nchandur/ranga/pull/1)]
* feat(board): hybrid board architecture combining parallel bitboards with a 64-square mailbox [[#1](https://github.com/nchandur/ranga/pull/1)]
