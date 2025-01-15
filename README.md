# Juice Benchmark

This repository contains benchmark tests comparing the performance of Juice ORM against GORM and standard database/sql operations.

## Benchmark Results

All tests run on Apple M1 (darwin/arm64)

### Single Record Creation Performance

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 9,100 | 133,040 | 535 | 15 |
| Juice | 10,000 | 140,911 | 2,937 | 61 |
| GORM | 10,000 | 144,790 | 4,215 | 52 |

### Batch Creation Performance (1000 Records)

| Framework | Operations | NS/Op     | B/Op      | Allocs/Op |
|-----------|------------|-----------|-----------|-----------|
| Standard DB | 211 | 6,041,218 | 578,805 | 35 |
| Juice | 175 | 7,620,343 | 1,334,215 | 21,838 |
| GORM | 130 | 8,931,485 | 1,494,792 | 13,062 |

### Query All Performance (1000 Records)

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 1,107 | 2,766,177 | 426,618 | 13,010 |
| Juice | 949 | 3,480,474 | 590,999 | 14,047 |
| GORM | 626 | 3,792,551 | 695,221 | 20,039 |

### Query With Limit Performance

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 1,327 | 2,654,317 | 336,534 | 8,673 |
| Juice | 1,185 | 3,052,607 | 501,862 | 9,722 |
| GORM | 637 | 4,269,654 | 557,873 | 20,037 |

### User Batch Create Performance

| Framework | Operations | NS/Op | B/Op | Allocs/Op |
|-----------|-----------|-------|------|-----------|
| Standard DB | 144 | 8,644,057 | 589,924 | 2,158 |
| Juice | 100 | 21,040,470 | 1,017,162 | 21,618 |
| GORM | 43 | 26,536,134 | 1,380,915 | 13,433 |

### Performance Visualization

```mermaid
graph TD
    subgraph Performance_Comparison_ns_op
        direction LR
        A[Single Create]
        B[Batch Create]
        C[Query All]
        D[Query Limit]
        E[User Batch]
    end

    classDef default fill:#f9f,stroke:#333,stroke-width:2px
    classDef performance fill:#dfd,stroke:#333,stroke-width:2px
    
    class A,B,C,D,E default
```

```mermaid
gantt
    title Performance Comparison (ns/op, lower is better)
    dateFormat X
    axisFormat %s

    section Single Create
    STD DB    : 0, 133040
    Juice     : 0, 140911
    GORM      : 0, 144790

    section Batch Create
    STD DB    : 0, 6041218
    Juice     : 0, 7620343
    GORM      : 0, 8931485

    section Query All
    STD DB    : 0, 2766177
    Juice     : 0, 3480474
    GORM      : 0, 3792551

    section Query Limit
    STD DB    : 0, 2654317
    Juice     : 0, 3052607
    GORM      : 0, 4269654

    section User Batch
    STD DB    : 0, 8644057
    Juice     : 0, 21040470
    GORM      : 0, 26536134
```

### Memory Usage and Allocations

| Operation | Framework | Memory (B/op) | Allocations (allocs/op) |
|-----------|-----------|---------------|------------------------|
| Single Create | STD DB | 535 | 15 |
|              | Juice  | 2,937 | 61 |
|              | GORM   | 4,215 | 52 |
| Batch Create | STD DB | 578,805 | 35 |
|              | Juice  | 1,334,215 | 21,838 |
|              | GORM   | 1,494,792 | 13,062 |
| Query All    | STD DB | 426,618 | 13,010 |
|              | Juice  | 590,999 | 14,047 |
|              | GORM   | 695,221 | 20,039 |
| Query Limit  | STD DB | 336,534 | 8,673 |
|              | Juice  | 501,862 | 9,722 |
|              | GORM   | 557,873 | 20,037 |
| User Batch   | STD DB | 589,924 | 2,158 |
|              | Juice  | 1,017,162 | 21,618 |
|              | GORM   | 1,380,915 | 13,433 |

### Analysis

#### Single Record Creation
- Standard DB performs best at ~133μs per operation
- Juice is about 5.9% slower than Standard DB
- GORM is about 8.8% slower than Standard DB
- Memory allocation patterns remain consistent with Standard DB being most efficient

#### Batch Creation (1000 Records)
- Standard DB leads at ~6.04ms per 1000 records
- Juice is about 26.1% slower than Standard DB
- GORM is about 47.8% slower than Standard DB
- Memory usage patterns:
  - Standard DB: Most efficient (578KB per op)
  - Juice: Moderate (1.33MB per op)
  - GORM: Highest (1.49MB per op)

#### Query Performance
- Standard DB processes queries in ~2.77ms per 1000 records
- Juice is about 25.8% slower than Standard DB
- GORM is about 37.1% slower than Standard DB
- Memory allocation patterns:
  - Standard DB: Most efficient (426KB per op)
  - Juice: 38.5% more memory than Standard DB
  - GORM: 63% more memory than Standard DB

#### Query With Limit Performance
- Standard DB shows best performance at ~2.65ms per operation
- Juice is about 15% slower than Standard DB
- GORM is about 60.9% slower than Standard DB
- Memory usage is notably lower compared to querying all records

#### User Batch Create Performance
- Standard DB shows best performance at ~8.64ms per operation
- Juice is about 143.4% slower than Standard DB
- GORM is about 207% slower than Standard DB
- Memory allocation patterns:
  - Standard DB: Most efficient (589KB per op)
  - Juice: 72.4% more memory than Standard DB
  - GORM: 134.1% more memory than Standard DB

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
