-- Origin SQL:
CREATE TABLE test.events_local UUID 'dad17568-b070-49d0-9ad1-7568b07029d0' (
    `date` Date,
    `f1` String,
    `f2` String,
    `f3` UInt64
    ) ENGINE = ReplacingMergeTree
    PARTITION BY date
    ORDER BY (f1, f2)
    SETTINGS index_granularity = 8192;

-- Format SQL:
CREATE TABLE test.events_local UUID 'dad17568-b070-49d0-9ad1-7568b07029d0' (`date` Date, `f1` String, `f2` String, `f3` UInt64) ENGINE = ReplacingMergeTree PARTITION BY date SETTINGS index_granularity=8192 ORDER BY (f1, f2);
