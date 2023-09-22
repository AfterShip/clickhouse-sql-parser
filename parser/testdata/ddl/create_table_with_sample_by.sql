CREATE TABLE default.test UUID '87887901-e33c-497e-8788-7901e33c997e'
(
    `f0` DateTime,
    `f1` UInt32,
    `f3` UInt32
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/{layer}/{shard}/default/test', '{replica}')
PARTITION BY toYYYYMM(timestamp)
ORDER BY (contractid, toDate(timestamp), userid)
SAMPLE BY userid
SETTINGS index_granularity = 8192;