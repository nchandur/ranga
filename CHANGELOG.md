# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
* feat(evaluate): static piece value evaluation
* feat(search): basic alpha-beta search
* feat(board): move generation
* feat(board): FEN parsing, `AddPiece`/`RemovePiece` utilities, and Zobrist hashing for position signatures.
* feat(board): pre-computed attack tables, including 1D arrays for leapers and magic bitboards for sliders
* feat(board): hybrid board architecture combining parallel Bitboards with a 64-square Mailbox.
