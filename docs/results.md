# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 8+0.08

**RC Engine:** ranga-1.20 | **Base Engine:** ranga-1.19

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=15, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 100.72 +/- 38.76                     |
| nElo          | 122.90 +/- 44.52                     |
| LOS           | 100.0%                               |
| Games         | 234 (W: 129, L: 63, D: 42)           |
| Score         | 150 / 234 (64.10%)                   |
| Draw ratio    | 35.90%                               |
| LLR           | 2.96 (100.6%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [6, 11, 42, 27, 31]                  |

Total time: 00:06:26 (h:m:s)

### LTC: 40+0.04

**RC Engine:** ranga-1.20 | **Base Engine:** ranga-1.19

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=15, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 85.94 +/- 34.47                      |
| nElo          | 109.23 +/- 41.91                     |
| LOS           | 100.0%                               |
| Games         | 264 (W: 130, L: 66, D: 68)           |
| Score         | 164 / 264 (62.12%)                   |
| Draw ratio    | 32.58%                               |
| LLR           | 3.01 (102.2%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [6, 17, 43, 39, 27]                  |

Total time: 00:08:36 (h:m:s)

## Progression Test vs *stash-21.0*
**Rounds:** 2500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 10+0.1 | **Threads:** 1 | **Hash:** 16 MB

| Version   | Elo               | nElo             | Wins | Losses | Draws | Score  |
|-----------|-------------------|------------------|------|--------|-------|--------|
| **v1.20** | 83.20 +/- 9.40    | 88.57 +/- 9.63   | 2842 | 1667   | 491   | 3087.5 |
| v1.19     | -78.58 +/- 8.91   | -87.88 +/- 9.63  | 1543 | 2655   | 802   | 1944   |
| v1.18     | -127.29 +/- 9.14  | -146.44 +/- 9.63 | 1262 | 3016   | 722   | 1623   |
| v1.17     | -165.71 +/- 9.60  | -192.64 +/- 9.63 | 1062 | 3281   | 657   | 1390.5 |
| v1.16     | -241.91 +/- 10.72 | -294.97 +/- 9.63 | 694  | 3704   | 602   | 995    |
