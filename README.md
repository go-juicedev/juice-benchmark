# Juice Benchmark

This repository contains benchmark tests comparing the performance of Juice ORM against GORM and standard database/sql operations.

## Benchmark Results

All tests run on Apple M1 (darwin/arm64)

### Single Record Creation Performance

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 8,991 | 144,952 | 535 | 15 |
| Juice | 10,000 | 147,548 | 2,272 | 59 |
| GORM | 10,000 | 148,690 | 4,215 | 52 |

### Batch Creation Performance (1000 Records)

| Framework | Operations | NS/Op     | B/Op      | Allocs/Op |
|-----------|------------|-----------|-----------|-----------|
| Standard DB | 202        | 6,158,115 | 578,805   | 35        |
| Juice | 144        | 7,713,170 | 1469,307  | 34,823    |
| GORM | 128        | 9,287,681 | 1,494,794 | 13,062    |

### Query Performance (1000 Records)

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 990 | 2,938,773 | 426,624 | 13,010 |
| Juice | 949 | 3,444,047 | 590,999 | 14,047 |
| GORM | 630 | 3,876,943 | 695,220 | 20,040 |

### Analysis

#### Single Record Creation
- All three frameworks show remarkably close performance (within ~3% of each other)
- Standard DB: Fastest at 144,952 ns/op with minimal memory usage (535 B/op)
- Juice: Only 1.8% slower than Standard DB
- GORM: 2.6% slower than Standard DB but with higher memory usage

#### Batch Creation (1000 Records)
- Standard DB performs significantly better: ~6.16ms per 1000 records
- Juice is about 36% slower than Standard DB but shows good throughput
- GORM is about 51% slower than Standard DB
- Memory usage varies significantly:
  - Standard DB: Most efficient (578KB per op)
  - GORM: Moderate (1.49MB per op)
  - Juice: Highest (2.28MB per op)

#### Query Performance
- Standard DB leads with 2.94ms per 1000-record query
- Juice is about 17% slower than Standard DB
- GORM is about 32% slower than Standard DB
- Memory allocation patterns:
  - Standard DB: Most efficient (416KB per op)
  - Juice: 38% more memory than Standard DB
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
