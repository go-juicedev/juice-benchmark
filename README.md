# Juice Benchmark

This repository contains benchmark tests comparing the performance of Juice ORM against GORM and standard database/sql operations.

## Benchmark Results

All tests run on Apple M1 (darwin/arm64)

### Single Record Creation Performance

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 8,398 | 145,237 | 535 | 15 |
| Juice | 10,000 | 144,127 | 2,938 | 61 |
| GORM | 10,000 | 156,914 | 4,215 | 52 |

### Batch Creation Performance (1000 Records)

| Framework | Operations | NS/Op     | B/Op      | Allocs/Op |
|-----------|------------|-----------|-----------|-----------|
| Standard DB | 205 | 6,343,141 | 578,807 | 35 |
| Juice | 157 | 7,482,669 | 1,334,188 | 21,838 |
| GORM | 124 | 9,133,590 | 1,494,789 | 13,062 |

### Query Performance (1000 Records)

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 1,136 | 3,044,905 | 426,618 | 13,010 |
| Juice | 948 | 3,177,366 | 590,994 | 14,047 |
| GORM | 621 | 4,099,117 | 695,196 | 20,039 |

### Analysis

#### Single Record Creation
- Juice now slightly outperforms Standard DB by about 0.8%
- Standard DB maintains minimal memory usage (535 B/op)
- GORM is about 8.8% slower than Juice
- Memory allocation patterns remain consistent

#### Batch Creation (1000 Records)
- Standard DB leads at ~6.34ms per 1000 records
- Juice is about 18% slower than Standard DB
- GORM is about 44% slower than Standard DB
- Memory usage patterns:
  - Standard DB: Most efficient (578KB per op)
  - Juice: Moderate (1.33MB per op)
  - GORM: Highest (1.49MB per op)

#### Query Performance
- Standard DB processes queries in ~3.04ms per 1000 records
- Juice is about 4.3% slower than Standard DB
- GORM is about 34.6% slower than Standard DB
- Memory allocation patterns:
  - Standard DB: Most efficient (426KB per op)
  - Juice: 38.5% more memory than Standard DB
  - GORM: 63% more memory than Standard DB

### Key Findings

1. **Single Record Operations**
   - All three solutions show excellent performance
   - Differences in speed are minimal (within 3%)
   - Memory usage varies significantly

2. **Batch Operations**
   - Standard DB shows clear advantages
   - Juice offers good performance with room for optimization
   - Memory usage could be optimized in both ORMs

3. **Query Performance**
   - All implementations handle 1000 records efficiently
   - Standard DB maintains the lead in both speed and memory
   - ORMs trade some performance for convenience

## Running the Benchmarks

To run the benchmarks:

```bash
go test -bench=. -benchmem
```

## Environment

- Go version: 1.21
- OS: Darwin/ARM64
- CPU: Apple M1
- MySQL: 8.0

> Note: Higher numbers in Operations indicate better performance, while lower numbers in NS/Op, B/Op, and Allocs/Op indicate better efficiency.
