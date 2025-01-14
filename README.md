# Juice Benchmark

This repository contains benchmark tests comparing the performance of Juice ORM against GORM and standard database/sql operations.

## Benchmark Results

All tests run on Apple M1 (darwin/arm64)

### Single Record Creation Performance

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 8,223 | 142,762 | 535 | 15 |
| Juice | 10,000 | 145,863 | 2,937 | 61 |
| GORM | 10,000 | 149,284 | 4,215 | 52 |

### Batch Creation Performance (1000 Records)

| Framework | Operations | NS/Op     | B/Op      | Allocs/Op |
|-----------|------------|-----------|-----------|-----------|
| Standard DB | 205 | 6,244,736 | 578,801 | 35 |
| Juice | 160 | 7,854,467 | 1,334,235 | 21,838 |
| GORM | 133 | 8,939,661 | 1,494,796 | 13,062 |

### Query Performance (1000 Records)

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 1,126 | 2,962,685 | 426,619 | 13,010 |
| Juice | 950 | 3,361,973 | 590,998 | 14,047 |
| GORM | 644 | 3,286,059 | 695,211 | 20,039 |

### Analysis

#### Single Record Creation
- All three frameworks show remarkably close performance (within ~4.5% of each other)
- Standard DB: Fastest at 142,762 ns/op with minimal memory usage (535 B/op)
- Juice: About 2.2% slower than Standard DB
- GORM: 4.6% slower than Standard DB but with higher memory usage

#### Batch Creation (1000 Records)
- Standard DB performs significantly better: ~6.24ms per 1000 records
- Juice is about 25.8% slower than Standard DB but shows good throughput
- GORM is about 43.2% slower than Standard DB
- Memory usage varies significantly:
  - Standard DB: Most efficient (578KB per op)
  - Juice: Moderate (1.33MB per op)
  - GORM: Highest (1.49MB per op)

#### Query Performance
- Standard DB leads with 2.96ms per 1000-record query
- Juice is about 13.5% slower than Standard DB
- GORM is about 10.9% slower than Standard DB
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
