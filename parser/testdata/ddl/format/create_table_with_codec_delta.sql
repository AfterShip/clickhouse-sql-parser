-- Origin SQL:
CREATE TABLE IF NOT EXISTS test_local
(
 `id` UInt64 CODEC(Delta, ZSTD(1)),
 `api_id` UInt64 CODEC(ZSTD(1)),
 `app_id` UInt64 CODEC(Delta(9), ZSTD(1)),
 `timestamp` DateTime64(9) CODEC(ZSTD(1)),
 INDEX timestamp_index(timestamp) TYPE minmax GRANULARITY 4
)
ENGINE = ReplicatedMergeTree('/root/test_local', '{replica}')
PARTITION BY toStartOfHour(`timestamp`)
ORDER BY (toUnixTimestamp64Nano(`timestamp`), `api_id`)
TTL toStartOfHour(`timestamp`) + INTERVAL 7 DAY,toStartOfHour(`timestamp`) + INTERVAL 2 DAY
SETTINGS execute_merges_on_single_replica_time_threshold=1200, index_granularity=16384, max_bytes_to_merge_at_max_space_in_pool=64424509440, storage_policy='main', ttl_only_drop_parts=1;


-- Format SQL:
CREATE TABLE IF NOT EXISTS test_local
(
  `id` UInt64 CODEC(Delta, ZSTD(1)),
  `api_id` UInt64 CODEC(ZSTD(1)),
  `app_id` UInt64 CODEC(Delta(9), ZSTD(1)),
  `timestamp` DateTime64(9) CODEC(ZSTD(1)),
  INDEX timestamp_index(timestamp) TYPE minmax GRANULARITY 4
)
ENGINE = ReplicatedMergeTree('/root/test_local', '{replica}')
PARTITION BY toStartOfHour(`timestamp`)
TTL toStartOfHour(`timestamp`) + INTERVAL 7 DAY,toStartOfHour(`timestamp`) + INTERVAL 2 DAY
SETTINGS execute_merges_on_single_replica_time_threshold=1200, index_granularity=16384, max_bytes_to_merge_at_max_space_in_pool=64424509440, storage_policy='main', ttl_only_drop_parts=1
ORDER BY (toUnixTimestamp64Nano(`timestamp`), `api_id`);
