-- Origin SQL:
CREATE TABLE test.`.inner.752391fb-44cc-4dd5-b523-91fb44cc9dd5`
    UUID '27673372-7973-44f5-a767-33727973c4f5' (
    `f0` String,
    `f1` String,
    `f2` LowCardinality(String),
    `f3` LowCardinality(String),
    `f4` DateTime64(3),
    `f5` Nullable(DateTime64(3)),
    `succeed_at` Nullable(DateTime64(3))
) ENGINE = MergeTree
PARTITION BY xxHash32(tag_id) % 20
ORDER BY label_id
SETTINGS index_granularity = 8192;


-- Format SQL:
CREATE TABLE test.`.inner.752391fb-44cc-4dd5-b523-91fb44cc9dd5` UUID '27673372-7973-44f5-a767-33727973c4f5' (`f0` String, `f1` String, `f2` LowCardinality(String), `f3` LowCardinality(String), `f4` DateTime64(3), `f5` Nullable(DateTime64(3)), `succeed_at` Nullable(DateTime64(3))) ENGINE = MergeTree PARTITION BY xxHash32(tag_id) % 20 SETTINGS index_granularity=8192 ORDER BY label_id;
