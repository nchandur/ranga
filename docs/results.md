# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 8+0.08

**RC Engine:** ranga-1.22 | **Base Engine:** ranga-1.21

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
| :------------ | -----------------------------------: |
| Result        | **H1 accepted** (pass)               |
| Elo           | 78.90 +/- 28.51                      |
| nElo          | 94.27 +/- 32.84                      |
| LOS           | 100.0%                               |
| Games         | 430 (W: 200, L: 104, D: 126)         |
| Score         | 263 / 430 (61.16%)                   |
| Draw ratio    | 32.56%                               |
| LLR           | 2.95 (100.3%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [14, 28, 70, 54, 49]                 |

Total time: 00:15:55 (h:m:s)

### LTC: 40+0.04

**RC Engine:** ranga-1.22 | **Base Engine:** ranga-1.21

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
| :------------ | -----------------------------------: |
| Result        | **H1 accepted** (pass)               |
| Elo           | 140.51 +/- 33.63                     |
| nElo          | 197.95 +/- 42.39                     |
| LOS           | 100.0%                               |
| Games         | 258 (W: 147, L: 48, D: 63)           |
| Score         | 178.5 / 258 (69.19%)                 |
| Draw ratio    | 29.46%                               | 
| LLR           | 2.97 (100.9%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [3, 7, 38, 50, 31]                   |

Total time: 00:31:05 (h:m:s)


## Progression Test vs *ranga-1.20*
**Rounds:** 500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 60+0.6 | **Threads:** 1 | **Hash:** 16 MB

| Version   | Elo               | nElo              | Wins | Losses | Draws | Score  |
|-----------|-------------------|-------------------|------|--------|-------|--------|
| v1.22     | 571.11 +/- 44.64  | 1142.56 +/- 21.53 | 936  |  8     |  56   |  964.0 |
| v1.21     | 502.57 +/- 38.80  | 879.69 +/- 21.53  | 914  |  19    |  67   |  947.5 |
