CREATE TABLE t0 on cluster default_cluster
(
    `method` Enum8('GET'=1 , 'POST'=2, 'HEAD'=3, 'PUT'=4,'PATCH'=5, 'DELETE'=6, 'CONNECT'=7, 'OPTIONS'=8, 'TRACE'=9) CODEC(ZSTD(1)),
    `timestamp` DateTime64(3) CODEC(DoubleDelta, ZSTD)
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}', '{replica}')
PARTITION BY toDate(timestamp)
ORDER BY (method,timestamp)
TTL toDate(timestamp) + toIntervalDay(3)
SETTINGS index_granularity = 8192;