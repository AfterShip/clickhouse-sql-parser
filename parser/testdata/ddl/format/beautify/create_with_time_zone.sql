-- Origin SQL:
CREATE TABLE IF NOT EXISTS test.db ON CLUSTER default_cluster
(
    `f0` Array(Tuple(
        f00 DateTime64(9, 'UTC'),
        f01 String,
        f02 Map(String, String),
        f03 Map(String, Float64),
        f04 Map(String, BOOL))) CODEC(ZSTD(1)
    ),
    `f1` UInt64 CODEC(Delta(8), LZ4),
    `f2` FixedString(16) CODEC(LZ4),
    `f3` FixedString(8) CODEC(LZ4),
    `f4` FixedString(8) CODEC(LZ4),
    `f6` DateTime64(9, 'UTC') CODEC(Delta(8), LZ4),
    `f6` UInt64 CODEC(Delta(8), LZ4),
    `f7` LowCardinality(String) CODEC(ZSTD(1)),
    `f8` String CODEC(ZSTD(1)),
    `f9` LowCardinality(String) CODEC(ZSTD(1)),
    `f10` String CODEC(ZSTD(1)),
    `f11` LowCardinality(String) CODEC(ZSTD(1)),
    `f12` LowCardinality(String) CODEC(ZSTD(1)),
    `f13` String CODEC(ZSTD(1)),
    `f14` Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    `f15` Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    `f16` Map(LowCardinality(String), Float64) CODEC(ZSTD(1)),
    `f17` Map(LowCardinality(String), BOOL) CODEC(ZSTD(1)),
    `f18` Array(Tuple(
        f180 FixedString(16),
        f181 FixedString(8),
        f182 String,
        f183 Map(String, String))) CODEC(ZSTD(1)),
        `f184` String CODEC(ZSTD(1)),
        `f185` String CODEC(ZSTD(1)),
        `f186` String CODEC(ZSTD(1)),
        `f187` UInt32 CODEC(ZSTD(1)),
        `f188` DATETIME DEFAULT now(),
        INDEX idx_0 f0 TYPE bloom_filter(0.001) GRANULARITY 1,
        INDEX idx_f1 f1 TYPE bloom_filter(0.001) GRANULARITY 1,
        INDEX idx_f2 f2 TYPE minmax GRANULARITY 1,
        INDEX idx_f3 f3 TYPE set(0) GRANULARITY 4,
        INDEX idx_f4 mapValues(f4) TYPE bloom_filter(0.01) GRANULARITY 1,
        INDEX idx_f5 name TYPE tokenbf_v1(4096, 3, 0) GRANULARITY 4
    )
    ENGINE = MergeTree
    PARTITION BY toDate(timestamp)
    ORDER BY (ts_bucket, service_name, name, toUnixTimestamp64Nano(timestamp))
    TTL toDate(timestamp) + toIntervalDay(15)
    SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1

-- Beautify SQL:
CREATE TABLE IF NOT EXISTS test.db ON CLUSTER default_cluster
(
  `f0` Array(Tuple(f00 DateTime64(9, 'UTC'), f01 String, f02 Map(String, String), f03 Map(String, Float64), f04 Map(String, BOOL))) CODEC(ZSTD(1)),
  `f1` UInt64 CODEC(Delta(8), LZ4),
  `f2` FixedString(16) CODEC(LZ4),
  `f3` FixedString(8) CODEC(LZ4),
  `f4` FixedString(8) CODEC(LZ4),
  `f6` DateTime64(9, 'UTC') CODEC(Delta(8), LZ4),
  `f6` UInt64 CODEC(Delta(8), LZ4),
  `f7` LowCardinality(String) CODEC(ZSTD(1)),
  `f8` String CODEC(ZSTD(1)),
  `f9` LowCardinality(String) CODEC(ZSTD(1)),
  `f10` String CODEC(ZSTD(1)),
  `f11` LowCardinality(String) CODEC(ZSTD(1)),
  `f12` LowCardinality(String) CODEC(ZSTD(1)),
  `f13` String CODEC(ZSTD(1)),
  `f14` Map(LowCardinality(String), String) CODEC(ZSTD(1)),
  `f15` Map(LowCardinality(String), String) CODEC(ZSTD(1)),
  `f16` Map(LowCardinality(String), Float64) CODEC(ZSTD(1)),
  `f17` Map(LowCardinality(String), BOOL) CODEC(ZSTD(1)),
  `f18` Array(Tuple(f180 FixedString(16), f181 FixedString(8), f182 String, f183 Map(String, String))) CODEC(ZSTD(1)),
  `f184` String CODEC(ZSTD(1)),
  `f185` String CODEC(ZSTD(1)),
  `f186` String CODEC(ZSTD(1)),
  `f187` UInt32 CODEC(ZSTD(1)),
  `f188` DATETIME DEFAULT now(),
  INDEX idx_0 f0 TYPE bloom_filter(0.001) GRANULARITY 1,
  INDEX idx_f1 f1 TYPE bloom_filter(0.001) GRANULARITY 1,
  INDEX idx_f2 f2 TYPE minmax GRANULARITY 1,
  INDEX idx_f3 f3 TYPE set(0) GRANULARITY 4,
  INDEX idx_f4 mapValues(f4) TYPE bloom_filter(0.01) GRANULARITY 1,
  INDEX idx_f5 name TYPE tokenbf_v1(4096, 3, 0) GRANULARITY 4
)
ENGINE = MergeTree
ORDER BY
  (ts_bucket, service_name, name, toUnixTimestamp64Nano(timestamp))
PARTITION BY toDate(timestamp)
TTL toDate(timestamp) + toIntervalDay(15)
SETTINGS
  index_granularity=8192,
  ttl_only_drop_parts=1;
