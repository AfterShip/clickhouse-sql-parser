ATTACH TABLE IF NOT EXISTS test.events_local ON CLUSTER 'default_cluster' (
    f0 String,
    f1 String,
    f2 String,
    f3 Datetime,
    f4 Datetime,
    f5 Map(String,String),
    f6 String,
    f7 Datetime DEFAULT now()
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}/test/events_local', '{replica}')
TTL f3 + INTERVAL 6 MONTH
PARTITION BY toYYYYMMDD(f3)
ORDER BY (f0,f1,f2);