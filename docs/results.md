# Results

## Contents
- [Ratings](#ratings)
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Ratings
1. `ranga-1.20`: [Link](https://cope-chess.live/engines/ranga/versions/1-20)
2. `ranga-1.22`: [Link](https://cope-chess.live/engines/ranga/versions/1-22)

## Sequential Probability Ratio Test

### STC: 8+0.08

**RC Engine:** ranga-1.23 | **Base Engine:** ranga-1.22

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=5, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
| :------------ | -----------------------------------: |
| Result        | **H1 accepted** (pass)               |
| Elo           | 48.12 +/- 15.61                      |
| nElo          | 58.81 +/- 18.83                      |
| LOS           | 100.0%                               |
| Games         | 1308 (W: 560, L: 380, D: 368)        |
| Score         | 744 / 1308 (56.88%)                  |
| Draw ratio    | 35.93%                               |
| LLR           | 2.95 (100.3%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [53, 96, 235, 158, 112]              |

Total time: 01:09:10 (h:m:s)

### LTC: 40+0.04

**RC Engine:** ranga-1.23 | **Base Engine:** ranga-1.22

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=5, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
| :------------ | -----------------------------------: |
| Result        | **H1 accepted** (pass)               |
| Elo           | 65.92 +/- 17.41                      |
| nElo          | 89.94 +/- 23.17                      |
| LOS           | 100.0%                               |
| Games         | 864 (W: 373, L: 211, D: 280)         |
| Score         | 513 / 864 (59.38%)                   |
| Draw ratio    | 37.27%                               | 
| LLR           | 2.97 (100.8%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [13, 69, 161, 121, 68]               |

Total time: 02:28:11 (h:m:s)

## Progression Test vs *ranga-1.22*
**Rounds:** 500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 60+0.6 | **Threads:** 1 | **Hash:** 16 MB

| Version   | Elo               | nElo              | Wins | Losses | Draws | Score  |
|-----------|-------------------|-------------------|------|--------|-------|--------|
| v1.23     | 87.26 +/- 16.30   | 120.28 +/- 21.53  | 449  |  203   |  348  |  623.0 |
