CREATE TABLE example1 (
    timestamp DateTime,
    x UInt32 TTL timestamp + INTERVAL 1 MONTH,
    y UInt32 TTL timestamp + INTERVAL 1 WEEK
)
ENGINE = MergeTree
ORDER BY tuple()