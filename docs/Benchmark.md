### Benchmark commands
```
./redis/src/redis-benchmark -p 3000 -t set -n 1000000 -r 1000000

./redis/src/redis-benchmark -n 1000000 -t get -c 500 -h localhost -p 3000 -r 1000000 --threads 3
```

### Benchmark result

Single threaded

```
Summary:
  throughput summary: 199163.52 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        1.401     0.136     1.191     2.303     6.047    44.383
```

Multi threaded
```
Summary:
  throughput summary: 116373.79 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        2.327     0.024     2.047     4.055    11.863    71.359
```

```
Summary:
  throughput summary: 128188.70 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        2.323     0.024     1.871     4.671     9.887    70.015
```

Sleep 100microsecond
- Multi threads

```
Summary:
  throughput summary: 26289.50 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
       18.994     4.112    18.383    19.743    36.735   169.983
```

- Single thread
```
Summary:
  throughput summary: 7044.78 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
       70.916     2.424    70.015    75.583    99.967   290.303
```
