# Results

## Contents
- [Elo Estimate](#elo-estimate)
- [Gauntlet](#round-robin)
- [SPRT](#sequential-probability-ratio-test)

---

## Elo Estimate

| Engine       | Reference Elo | Estimated Elo    |
|--------------|---------------|------------------|
| Stash-8.0    | 1090          | 1090             |
| **ranga-1.0**| —637          | **~453 +/- 36**  |
| ranga-1.1    |               | **~473 +/- 36**  | 
> **Estimated Elo: ~417-489**

---

## Round Robin

**Date:** 20/08/2026
 
**Engines:** ranga-1.0, ranga-1.1, stash-8.0
**TC:** 8+0.08 · **Rounds:** 1000 (games/encounter=2, concurrency=15)
 
| Rank | Engine     | Elo     | +/-   | nElo    | +/-   | Games | Score | Draw  |
|------|------------|---------|-------|---------|-------|-------|-------|-------|
| 1    | stash-8.0  | 964.02  | 63.58 | 3946.72 | 10.77 | 4000  | 99.6% | 0.00% |
| 2    | ranga-1.1  | -182.16 | 10.62 | -220.62 | 10.77 | 4000  | 25.9% | 41.2% |
| 3    | ranga-1.0  | -196.10 | 10.43 | -248.36 | 10.77 | 4000  | 24.4% | 41.2% |
 
**Ptnml (0-2):**
 
| Engine     | LL  | LD  | DD  | WD  | WW   |
|------------|-----|-----|-----|-----|------|
| stash-8.0  |   0 |   0 |   0 |  31 | 1969 |
| ranga-1.1  | 987 |  71 | 824 | 115 |    3 |
| ranga-1.0  | 986 | 132 | 824 |  57 |    1 |
 
Total time: 01:15:31 (h:m:s)

---

## Sequential Probability Ratio Test

**Date:** 20/08/2026

**RC Engine:** ranga-1.1 | **Base Engine:** ranga-1.0

**TC:** 8+0.08 | **Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=5, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 19.94 +/- 5.65                       |
| nElo          | 68.00 +/- 19.21                      |
| LOS           | 100.00%                              |
| Games         | 1256 (W: 94, L: 22, D: 1140)         |
| Score         | 664 / 1256 (52.87%)                  |
| Draw ratio    | 82.01%                               |
| LLR           | 2.96 (100.4%) — bounds (-2.94, 2.94) |

Total time: 00:29:15 (h:m:s)
