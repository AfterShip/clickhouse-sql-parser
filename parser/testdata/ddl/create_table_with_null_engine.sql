CREATE TABLE logs.t0 on cluster default
(
    `trace_id` String CODEC(ZSTD(1)),
    INDEX trace_id_bloom_idx trace_id TYPE bloom_filter(0.01) GRANULARITY 64
) ENGINE = Null();