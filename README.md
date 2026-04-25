# Godis

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Redis Protocol](https://img.shields.io/badge/RESP-compatible-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/docs/latest/develop/reference/protocol-spec/)
[![Architecture](https://img.shields.io/badge/Architecture-shared--nothing-black?style=for-the-badge)](#architecture)

Godis is a Redis-like server written from scratch in Go. It implements RESP command parsing, in-memory Redis data structures, TTL handling, eviction experiments, probabilistic data structures, and an evolving multi-threaded shared-nothing runtime.

The goal is not to wrap Redis. The goal is to build the core pieces of a Redis-style database server by hand and make the tradeoffs visible.

## Highlights

- RESP request parsing and Redis CLI compatibility
- Multi-listener server path with I/O handlers and worker shards
- Shared-nothing command execution: each worker owns an independent `RedisDB`
- Platform I/O multiplexing wrappers for Linux `epoll` and macOS `kqueue`
- Strings, sets, sorted sets, Bloom filters, and Count-Min Sketch commands
- TTL commands and per-database expiration support
- Eviction policy experiments, including LRU sampling
- Benchmark and profiling notes under `docs/`

## Architecture

```text
redis-cli / redis-benchmark
          |
          v
   TCP listeners (:3000)
          |
          v
     I/O handlers
          |
          v
  key hash dispatcher
          |
          v
  worker 0   worker 1   worker N
  RedisDB    RedisDB    RedisDB
```

The active server path is designed around a shared-nothing shard model:

1. Listeners accept client connections on port `3000`.
2. I/O handlers read RESP commands from sockets.
3. Commands are dispatched by hashing the first key argument.
4. Each worker executes commands serially against its own `RedisDB` shard.

This keeps `RedisDB` simple: it is owned by one worker and does not need internal locking for normal command execution.

## Supported Commands

| Category | Commands |
| --- | --- |
| Core | `PING`, `INFO` |
| Strings | `SET`, `GET`, `DEL`, `EXISTS` |
| Expiration | `EXPIRE`, `TTL`, `PTTL` |
| Sets | `SADD`, `SREM`, `SISMEMBER`, `SMEMBERS` |
| Sorted sets | `ZADD`, `ZSCORE`, `ZRANK`, `ZREM` |
| Bloom filter | `BF.RESERVE`, `BF.ADD`, `BF.MADD`, `BF.EXISTS`, `BF.MEXISTS` |
| Count-Min Sketch | `CMS.INITBYDIM`, `CMS.INITBYPROB`, `CMS.INCRBY`, `CMS.QUERY` |

Note: multi-key commands are still evolving in the sharded runtime. Single-key commands route to the owning worker. Cross-shard multi-key semantics need explicit coordination before they can be considered Redis-compatible.

## Quick Start

### Requirements

- Go 1.25+
- macOS or Linux
- Optional: Redis CLI or Redis benchmark tools

### Run the server

```sh
go run ./cmd
```

The server listens on:

```text
localhost:3000
```

The pprof endpoint is also enabled while the server is running:

```text
localhost:6060/debug/pprof
```

### Connect with Redis CLI

```sh
redis-cli -p 3000
```

Example session:

```redis
127.0.0.1:3000> PING
PONG
127.0.0.1:3000> SET user:1 Ada
OK
127.0.0.1:3000> GET user:1
"Ada"
127.0.0.1:3000> EXPIRE user:1 10
(integer) 1
127.0.0.1:3000> TTL user:1
(integer) 9
```

## Development

Run the test suite with a writable Go cache:

```sh
GOCACHE=/tmp/godis-gocache go test ./...
```

Run benchmarks against the server:

```sh
redis-benchmark -n 1000000 -t get,set -c 500 -h localhost -p 3000 -r 1000000 --threads 3
```

More notes:

- [Benchmark notes](docs/Benchmark.md)
- [Profiling notes](docs/Profiling.md)
- [Redis CLI setup](docs/Redis_CLI.md)

## Project Layout

```text
.
|-- cmd/                         # Server entrypoint
|-- internal/
|   |-- config/                  # Runtime constants
|   |-- constant/                # Command and server constants
|   |-- core/                    # RESP, executor, RedisDB, commands, workers
|   |   |-- data_structure/      # Dict, skiplist, sorted set, Bloom, CMS, eviction
|   |   `-- io_multiplexer/      # epoll/kqueue abstraction
|   `-- server/                  # Listeners, I/O handlers, shutdown flow
|-- docs/                        # Benchmarks, profiling, CLI notes
|-- Signal/                      # Historical experiment
|-- ThreadPerConn/               # Historical experiment
`-- ThreadPool/                  # Historical experiment
```

## Design Notes

Godis is intentionally built as a learning-focused systems project. The code favors explicit ownership boundaries over hidden synchronization:

- Workers own their database shards.
- I/O handlers own socket readiness and connection reads.
- Dispatch is based on key hashing.
- Expiration and eviction should remain local to the owning shard.
- Shared mutable state across workers is avoided unless there is a clear coordination design.

## Current Focus

The project is moving from a single-threaded event loop toward a multi-threaded shared-nothing architecture. The most important correctness boundary is preserving Redis-like command behavior while routing work across independent shards.

Areas still worth improving:

- Cross-shard semantics for multi-key commands
- Async worker replies with per-connection response ordering
- More command coverage
- Broader integration tests through `redis-cli` and `redis-benchmark`
- Clearer runtime configuration instead of compile-time constants

## License

This repository does not currently declare a license. Add one before publishing or distributing derived work.
