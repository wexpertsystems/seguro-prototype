# Dqlite

[dqlite](https://dqlite.io) is a distributed Sqlite technology 
published by Canonical. 

# Dependencies

1. Linux (Ubuntu is easiest, others must build dqlite from source)
2. [go 1.18](https://go.dev/dl/)
3. [dqlite and its dependencies](https://github.com/canonical/dqlite) (build from source or install via development PPA on Ubuntu)
4. Execution on the host network (containerization/virtualization possible, but tedious)

# Build and Run Instructions

```
export CGO_LDFLAGS_ALLOW="-Wl,-z,now"
go build -tags libsqlite3 -o seguro-dqlite.out
./seguro-dqlite.out
```