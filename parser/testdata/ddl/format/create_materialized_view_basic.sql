-- Origin SQL:
ATTACH MATERIALIZED VIEW test.events_local
UUID '3493e374-e2bb-481b-b493-e374e2bb981b'
(`f0` DateTime64(3),
`f1` String,
`f2` String,
`f3` String,
`f4` String,
`f5` Int64)
ENGINE = ReplicatedAggregatingMergeTree('/clickhouse/tables/{layer}-{shard}}')
PARTITION BY toDate(f1)
ORDER BY (f1, f2, f3, f4)
SETTINGS index_granularity = 8192;

-- Format SQL:
CREATE MATERIALIZED VIEW test.events_local
(
    `f0` DateTime64(3),
    `f1` String,
    `f2` String,
    `f3` String,
    `f4` String,
    `f5` Int64
)
ENGINE = ReplicatedAggregatingMergeTree('/clickhouse/tables/{layer}-{shard}}')
PARTITION BY toDate(f1)
SETTINGS index_granularity=8192
ORDER BY (f1, f2, f3, f4);
