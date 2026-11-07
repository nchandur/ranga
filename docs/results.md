# Results

## Contents
- [SPRT](#sequential-probability-ratio-test)
- [Progression Test](#progression-test)

## Sequential Probability Ratio Test

### STC: 10+0.1

**RC Engine:** ranga-1.21 | **Base Engine:** ranga-1.20

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
| :------------ | -----------------------------------: |
| Result        | **H1 accepted** (pass)               |
| Elo           | 173.75 +/- 43.15                     |
| nElo          | 210.64 +/- 44.14                     |
| LOS           | 100.0%                               |
| Games         | 238 (W: 152, L: 42, D: 44)           |
| Score         | 174 / 238 (73.11%)                   |
| Draw ratio    | 26.89%                               |
| LLR           | 2.98 (101.1%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [4, 6, 32, 30, 47]                   |

Total time: 00:05:33 (h:m:s)

### LTC: 40+0.4

**RC Engine:** ranga-1.21 | **Base Engine:** ranga-1.20

**Rounds:** 10000 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**SPRT:** elo0=0, elo1=10, $\alpha$=0.05, $\beta$=0.05


| Metric        | Value                                |
| :------------ | -----------------------------------: |
| Result        | **H1 accepted** (pass)               |
| Elo           | 470.44 +/- 84.24                     |
| nElo          | 882.20 +/- 53.83                     |
| LOS           | 100.0%                               |
| Games         | 160 (W: 142, L: 2, D: 16)            |
| Score         | 150 / 160 (93.75%)                   |
| Draw ratio    | 2.50%                                | 
| LLR           | 2.98 (101.2%) — bounds (-2.94, 2.94) |
| Ptnml(0-2)    | [0, 0, 2, 16, 62]                    |

Total time: 00:10:34 (h:m:s)

## Progression Test vs *ranga-1.20*
**Rounds:** 500 (games=2) | **Book:** `UHO_Lichess_4852_v1.epd`

**TC:** 60+0.6 | **Threads:** 1 | **Hash:** 16 MB

| Version   | Elo               | nElo             | Wins | Losses | Draws | Score  |
|-----------|-------------------|------------------|------|--------|-------|--------|
| v1.21     | 502.57 +/- 38.80  | 879.69 +/- 21.53 | 914  |  19    |  67   |  947.5 |
