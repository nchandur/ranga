# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 8+0.08

**RC Engine:** ranga-1.18 | **Base Engine:** ranga-1.17

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 14.82 +/- 10.44                      |
| nElo          | 19.42 +/- 13.66                      |
| LOS           | 99.73%                               |
| Games         | 2486 (W: 903, L: 797, D: 786)        |
| Score         | 1296 / 2486 (52.13%)                 |
| Draw ratio    | 39.34%                               |
| LLR           | 2.96 (100.5%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [99, 243, 489, 277, 135]             |

Total time: 00:55:11 (h:m:s)

### LTC: 40+0.04

**RC Engine:** ranga-1.18 | **Base Engine:** ranga-1.17

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 12.11 +/- 8.92                       |
| nElo          | 16.78 +/- 12.35                      |
| LOS           | 99.61%                               |
| Games         | 3042 (W: 1075, L: 969, D: 998)       |
| Score         | 1574 / 3042 (51.74%)                 |
| Draw ratio    | 43.59%                               |
| LLR           | 2.96 (100.5%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [110, 278, 663, 336, 134]            |

Total time: 02:00:49 (h:m:s)

## Progression Test vs *stash-21.0*
**Rounds:** 2500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 10+0.1 | **Threads:** 1 | **Hash:** 16 MB

| Metric | *v1.16*                  | *v1.17*                    | *v1.18*                    |
| ------ | ------------------------ | -------------------------- | -------------------------- |
| Elo    | -241.91 +/- 10.72        | -165.71 +/- 9.60           | -127.29 +/- 9.14           |
| nElo   | -294.97 +/- 9.63         | -192.64 +/- 9.63           | -146.44 +/- 9.63           |
| LOS    | 0.00                     | 0.00                       | 0.00                       |
| Wins   | 694                      | 1062                       | 1262                       |
| Losses | 3704                     | 3281                       | 3016                       |
| Draws  | 602                      | 657                        | 722                        |
| Score  | 995                      | 1390.5                     | 1623                       |
| WL/DD  | 17.63                    | 18.35                      | 14.03                      |
| Ptnml  | [1358, 459, 559, 83, 41] | [1047, 453, 774, 124, 102] | [879, 444, 872, 162, 143]  |
