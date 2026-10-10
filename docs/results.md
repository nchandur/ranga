# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 8+0.08

**Date:** 2026-10-10

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

**Date:** 2026-10-10

**RC Engine:** ranga-1.17 | **Base Engine:** ranga-1.16

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
|---------------|--------------------------------------|
| Result        | **H1 accepted** (pass)               |
| Elo           | 39.89 +/- 19.25                      |
| nElo          | 50.51 +/- 24.14                      |
| LOS           | 100.0%                               |
| Games         | 796 (W: 318, L: 227, D: 251)         |
| Score         | 443.5 / 796 (55.72%)                 |
| Draw ratio    | 35.93%                               |
| LLR           | 2.97 (101.0%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [23, 79, 143, 90, 63]                |

Total time: 01:08:37 (h:m:s)

## Progression Test

### ranga-1.16 vs stash-21.0

**Date:** 2026-10-08

**Rounds:** 2500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 10+0.01 | **Threads:** 1 | **Hash:** 16 MB

| Metric             | Value                                |
|--------------------|--------------------------------------|
| Elo                | -241.91 +/- 10.72                    |
| nElo               | -294.97 +/- 9.63                     |
| LOS                | 0.00%                                |
| Games              | 5000 (W: 694, L: 3704, D: 602)       |
| Score              | 995 / 5000 (19.90%)                  |
| Draw ratio (pairs) | 22.36%                               |
| WL/DD              | 17.63                                |
| Ptnml(0-2)         | [1358, 459, 559, 83, 41]             |

Total time: 01:24:35 (h:m:s)
