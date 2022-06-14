# Dqlite

[dqlite](https://dqlite.io) is a distributed Sqlite technology 
published by Canonical. 

# Dependencies

1. Linux (Ubuntu is easiest, others must build dqlite from source)
2. [go 1.18](https://go.dev/dl/)
3. [dqlite and its dependencies](https://github.com/canonical/dqlite) (build from source or install via development PPA on Ubuntu)
4. Execution on the host network (containerization/virtualization possible, but tedious)

# Instructions

## Build

```
export CGO_LDFLAGS_ALLOW="-Wl,-z,now"
go build -tags libsqlite3 ./cmd/dqlite-benchmark
```

## Run

Single-node benchmark:
```
./dqlite-benchmark -d 127.0.0.1:9001 --driver --cluster 127.0.0.1:9001
```

Run a multi-node benchmark with the master as the driver:
```
./dqlite-benchmark --db 127.0.0.1:9001 --driver --cluster 127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003 &
./dqlite-benchmark --db 127.0.0.1:9002 --join 127.0.0.1:9001 &
./dqlite-benchmark --db 127.0.0.1:9003 --join 127.0.0.1:9001 &
```

Run a multi-node benchmark with a replica node as the driver:
```
./dqlite-benchmark --db 127.0.0.1:9001 &
./dqlite-benchmark --db 127.0.0.1:9002 --join 127.0.0.1:9001 &
./dqlite-benchmark --db 127.0.0.1:9003 --join 127.0.0.1:9001 --driver --cluster 127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003 &
```

The results can be found on the `driver` node in `/tmp/dqlite-benchmark/127.0.0.1:9001/results` or in the directory provided to the tool.
Benchmark results are files named `n-q-timestamp` where `n` is the number of the worker,
`q` is the type of query that was tracked. All results in the file are in milliseconds.