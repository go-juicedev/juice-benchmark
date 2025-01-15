# Juice Benchmark

This repository contains benchmark tests comparing the performance of Juice ORM against GORM and standard database/sql operations.

## Benchmark Results

All tests run on Apple M1 (darwin/arm64)

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

```mermaid
gantt
    title Memory Usage Comparison (B/op, lower is better)
    dateFormat X
    axisFormat %s

    section Single Create
    STD DB    : 0, 535
    Juice     : 0, 2937
    GORM      : 0, 4215

    section Batch Create
    STD DB    : 0, 578805
    Juice     : 0, 1334215
    GORM      : 0, 1494792

    section Query All
    STD DB    : 0, 426618
    Juice     : 0, 590999
    GORM      : 0, 695221

    section Query Limit
    STD DB    : 0, 336534
    Juice     : 0, 501862
    GORM      : 0, 557873

    section User Batch
    STD DB    : 0, 589924
    Juice     : 0, 1017162
    GORM      : 0, 1380915
```

```mermaid
gantt
    title Memory Allocations Comparison (allocs/op, lower is better)
    dateFormat X
    axisFormat %s

    section Single Create
    STD DB    : 0, 15
    Juice     : 0, 61
    GORM      : 0, 52

    section Batch Create
    STD DB    : 0, 35
    Juice     : 0, 21838
    GORM      : 0, 13062

    section Query All
    STD DB    : 0, 13010
    Juice     : 0, 14047
    GORM      : 0, 20039

    section Query Limit
    STD DB    : 0, 8673
    Juice     : 0, 9722
    GORM      : 0, 20037

    section User Batch
    STD DB    : 0, 2158
    Juice     : 0, 21618
    GORM      : 0, 13433
```

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
