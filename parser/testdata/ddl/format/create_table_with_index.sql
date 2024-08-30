-- Origin SQL:
CREATE TABLE IF NOT EXISTS test_local
(
 `id` UInt64 CODEC(Delta, ZSTD(1)),
 `api_id` UInt64 CODEC(ZSTD(1)),
 `arr` Array(Int64),
 `content` String CODEC(ZSTD(1)),
 `output` String,
 INDEX id_idx id TYPE minmax GRANULARITY 10,
 INDEX api_id_idx api_id TYPE set(100) GRANULARITY 2,
 INDEX arr_idx arr TYPE bloom_filter(0.01) GRANULARITY 3,
 INDEX content_idx content TYPE tokenbf_v1(30720, 2, 0) GRANULARITY 1,
 INDEX output_idx output TYPE ngrambf_v1(3, 10000, 2, 1) GRANULARITY 2
)
ENGINE = ReplicatedMergeTree('/root/test_local', '{replica}')
PARTITION BY toStartOfHour(`timestamp`)
ORDER BY (toUnixTimestamp64Nano(`timestamp`), `api_id`)
TTL toStartOfHour(`timestamp`) + INTERVAL 7 DAY,toStartOfHour(`timestamp`) + INTERVAL 2 DAY
SETTINGS execute_merges_on_single_replica_time_threshold=1200, index_granularity=16384, max_bytes_to_merge_at_max_space_in_pool=64424509440, storage_policy='main', ttl_only_drop_parts=1;


-- Format SQL:
CREATE TABLE IF NOT EXISTS test_local (`id` UInt64 CODEC(Delta, ZSTD(1)), `api_id` UInt64 CODEC(ZSTD(1)), `arr` Array(Int64), `content` String CODEC(ZSTD(1)), `output` String, INDEX id_idx id TYPE minmax GRANULARITY 10, INDEX api_id_idx api_id TYPE set(100) GRANULARITY 2, INDEX arr_idx arr TYPE bloom_filter(0.01) GRANULARITY 3, INDEX content_idx content TYPE tokenbf_v1(30720, 2, 0) GRANULARITY 1, INDEX output_idx output TYPE ngrambf_v1(3, 10000, 2, 1) GRANULARITY 2) ENGINE = ReplicatedMergeTree('/root/test_local', '{replica}') PARTITION BY toStartOfHour(`timestamp`) TTL toStartOfHour(`timestamp`) + INTERVAL 7 DAY, toStartOfHour(`timestamp`) + INTERVAL 2 DAY SETTINGS execute_merges_on_single_replica_time_threshold=1200, index_granularity=16384, max_bytes_to_merge_at_max_space_in_pool=64424509440, storage_policy='main', ttl_only_drop_parts=1 ORDER BY (toUnixTimestamp64Nano(`timestamp`), `api_id`);
