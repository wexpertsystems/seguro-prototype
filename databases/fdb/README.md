# FoundationDB Benchmark

FoundationDB is a distributed, multi-model data store by Apple.

# Dependencies

1. A local N-node FoundationDB cluster.
   [Installation instructions here.](https://www.foundationdb.org/)
   [3-node configuration here.](https://apple.github.io/foundationdb/configuration.html?highlight=configuration#fdbserver-id-section-s)
   a. Your `foundationdb.conf` file should have N `[fdbserver.<ID>]` sections,
   depending on the type of benchmark you want to run. sections, one for each
   node. If you change the `.conf` file, the FoundationDB service must be
   restarted for changes to take effect.
2. [go 1.18](https://go.dev/dl/)
3. Execution on the host network

# Instructions

## Build

```
go build -o fdb-benchmark
```

## Run

```
For benchmarking FoundationDB.

Usage:
  fdb-benchmark [flags]

Flags:
      --batch-size int    Number of events per batch. (default 5)
      --events int        Number of events to write to the database. (default 4096)
  -h, --help              help for fdb-benchmark
      --value-size int    Size of the mock event values in bytes. (default 1024)
  -w, --workload string   The workload to run: "single", or "batch". (default "single")
```
