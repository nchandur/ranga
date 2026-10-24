# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 8+0.08

**RC Engine:** ranga-1.19 | **Base Engine:** ranga-1.18

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=15, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 47.20 +/- 25.68                      |
| nElo          | 58.32 +/- 31.28                      |
| LOS           | 99.99%                               |
| Games         | 474 (W: 213, L: 149, D: 112)         |
| Score         | 269 / 474 (56.75%)                   |
| Draw ratio    | 40.51%                               |
| LLR           | 2.96 (100.4%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [15, 39, 96, 41, 46]                 |

Total time: 00:09:13 (h:m:s)

### LTC: 40+0.04

**RC Engine:** ranga-1.19 | **Base Engine:** ranga-1.18

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=15, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 40.70 +/- 23.59                      |
| nElo          | 47.47 +/ 27.22                       |
| LOS           | 99.97%                               |
| Games         | 626 (W: 272, L: 199, D: 155)         |
| Score         | 349.5 / 626 (55.83%)                 |
| Draw ratio    | 32.27%                               |
| LLR           | 3.03 (102.8%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [33, 47, 101, 78, 54]                |

Total time: 00:12:38 (h:m:s)

## Progression Test vs *stash-21.0*
**Rounds:** 2500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 10+0.1 | **Threads:** 1 | **Hash:** 16 MB

| Version   | Elo               | nElo             | Wins | Losses | Draws | Score  |
|-----------|-------------------|------------------|------|--------|-------|--------|
| **v1.19** | -78.58 +/- 8.91   | -87.88 +/- 9.63  | 1543 | 2655   | 802   | 1944   |
| v1.18     | -127.29 +/- 9.14  | -146.44 +/- 9.63 | 1262 | 3016   | 722   | 1623   |
| v1.17     | -165.71 +/- 9.60  | -192.64 +/- 9.63 | 1062 | 3281   | 657   | 1390.5 |
| v1.16     | -241.91 +/- 10.72 | -294.97 +/- 9.63 | 694  | 3704   | 602   | 995    |
