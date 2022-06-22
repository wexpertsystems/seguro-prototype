# Seguro Prototype

[Prototype Grant](https://urbit.org/grants/seguro-prototype)  
[Seguro & Armada Whitepaper](https://gist.github.com/wexpert/0485a722185d5ee70742570036faf32f)

## Overview

This repository implements benchmarks for [Dqlite](https://dqlite.io) and
[FoundationDB](https://foundationdb.org) to determine their suitability for
integration with `vere` as part of the Seguro project.

## Databases

### Dqlite

Dqlite is a distributed version of SQLite developed by Canonical. It adds a
networking layer above SQLite to achieve data replication while retaining
SQLite's embedded nature and also its ACID transaction compliance.

### FoundationDB

FoundationDB is a distributed, multi-model database developed by Apple. It
maintains ACID compliance and has a reputation for impressive performance and
scalability.

## Benchmarks

For the performance benchmarks, the test data was designed to simulate Urbit
events as such:

```
<continuous integer keys> -> <arbitrarily sized byte buffer values>
```

Dqlite currently only has an official Go client available, with a C client
planned for release with Ubuntu 20.10. FoundationDB has an official Go client as
well. In order to control as many variables as possible, both benchmarks were
implemented in Go. Instructions for building and execution can be found in the
respective database directories.

With this simple data model, which database exhibits the best overall
performance?

Note that we only care about write performance because database reads are not
involved in the majority of Urbit's activities. Events are processed as follows:

1. Compute event
2. Commit event
3. Release effect

All benchmarks were performed on the same
[Linux box](https://www.amazon.com/Windows-Desktop-Computer-2500Mbps-Graphics/dp/B093V18HKB):

- AMD Ryzen 5 4500U (6C/6T up to 4.0GHz)
- 2x8GB DDR4 RAM
- Samsung 980 PRO NVMe SSD
- Ubuntu 20.04

### Single writes

| database | dqlite    | fdb       |
| -------- | --------- | --------- |
| events   | 4301      | 4096      |
| errors   | 0         | 0         |
| avg [ms] | 13.938089 | 6.344161  |
| max [ms] | 91.736463 | 16.663706 |
| min [ms] | 8.023056  | 2.627436  |

Winner: FoundationDB

### Batch writes

`vere` often batches multiple events into single database commits when it can,
so the results here are of particular importance.

| database   | dqlite    | fdb      |
| ---------- | --------- | -------- |
| events     | 4208      | 4096     |
| batch size | 5         | 5        |
| errors     | 0         | 0        |
| avg [ms]   | 14.229881 | 0.739242 |
| max [ms]   | 85.146681 | 2.399221 |
| min [ms]   | 8.154588  | 0.573965 |

Winner: FoundationDB

### Fragmented single writes

FoundationDB supports a maximum value size of 100KB and a recommended size of <
10KB. For events that exceed these limits, we use a fragmentation algorithm to
store them across multiple underlying database values. These benchmarks use a
maximum database value size of 10KB, so a 100KB event would be split across 10
or so underlying values. Dqlite value sizes are practically unlimited, so we
don't include fragmented benchmarks for it.

| database       | fdb      |
| -------------- | -------- |
| events         | 4096     |
| event size     | 100KB    |
| max value size | 10KB     |
| errors         | 0        |
| avg [ms]       | 4.816258 |
| max [ms]       | 3.776849 |
| min [ms]       | 0.572935 |

### Fragmented batch writes

| database       | fdb      |
| -------------- | -------- |
| events         | 4096     |
| batch size     | 5        |
| event size     | 100KB    |
| max value size | 10KB     |
| errors         | 0        |
| avg [ms]       | 0.689704 |
| max [ms]       | 1.447475 |
| min [ms]       | 0.574091 |

## Conclusion

**For Seguro, the right choice is FoundationDB.**

Even with Dqlite's embedded architecture and lack of need for a fragmentation
scheme, FoundationDB far outpaces it in both single and batch event writes.
FoundationDB is 4x as fast on single writes, and 16x as fast with batches.

In the minds of various Urbit runtime engineers over the years, Dqlite has long
been considered the distributed database of choice for integration in `vere` for
implementation of a fault-tolerant Urbit. These results indicate otherwise,
however.

In fact, FoundationDB is preferable in all of the following regards:

- Performance (see above)
- Production-readiness (FoundationDB is used in Apple CloudKit; Dqlite is a new,
  unproven technology)
- Documentation (Apple; Canonical)
- API client language variety (FoundationDB has C/Go/Python/Ruby/Java; Dqlite
  only has Go (C planned for future release))
