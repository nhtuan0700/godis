package config

const Protocol = "tcp"
const Address = ":3000"
const MaxConnections = 20000

const EvictionRatio = 0.1
const MaxKeyNumber = 10

// policy: "allkeys-random" | "allkeys-lru"
const EvictPolicy = "allkeys-lru"

const EpoolMaxSize = 16
const LruSampledSize = 5
