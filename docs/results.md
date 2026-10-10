# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 8+0.08

**RC Engine:** ranga-1.17 | **Base Engine:** ranga-1.16

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 67.23 +/- 26.91                      |
| nElo          | 76.81 +/- 29.92                      |
| LOS           | 100.0%                               |
| Games         | 518 (W: 254, L: 155, D: 109)         |
| Score         | 308.5 / 518 (59.56%)                 |
| Draw ratio    | 35.52%                               |
| LLR           | 3.00 (101.8%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [18, 41, 92, 40, 68]                 |

Total time: 00:17:39 (h:m:s)

### LTC: 40+0.04

**RC Engine:** ranga-1.17 | **Base Engine:** ranga-1.16

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 45.41 +/- 21.27                      |
| nElo          | 53.63 +/- 24.80                      |
| LOS           | 100.0%                               |
| Games         | 754 (W: 334, L: 236, D: 184)         |
| Score         | 426 / 754 (56.50%)                   |
| Draw ratio    | 35.01%                               |
| LLR           | 2.95 (100.3%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [36, 54, 132, 86, 69]                |

Total time: 00:41:42 (h:m:s)

## Progression Test vs *stash-21.0*
**Rounds:** 2500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 10+0.1 | **Threads:** 1 | **Hash:** 16 MB

| Metric | *v1.16*                  | *v1.17*                    |
| ------ | ------------------------ | -------------------------- |
| Elo    | -241.91 +/- 10.72        | -165.71 +/- 9.60           |
| nElo   | -294.97 +/- 9.63         | -192.64 +/- 9.63           |
| LOS    | 0.00                     | 0.00                       |
| Wins   | 694                      | 1062                       |
| Losses | 3704                     | 3281                       |
| Draws  | 602                      | 657                        |
| Score  | 995                      | 1390.5                     |
| WL/DD  | 17.63                    | 18.35                      |
| Ptnml  | [1358, 459, 559, 83, 41] | [1047, 453, 774, 124, 102] |
