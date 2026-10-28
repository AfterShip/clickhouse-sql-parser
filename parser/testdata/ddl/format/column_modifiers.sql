-- Origin SQL:
CREATE TABLE db.example
(
    event_time DateTime,
    value_a String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY,
    value_b String COMMENT 'example' TTL event_time + INTERVAL 1 DAY CODEC(ZSTD(1)),
    value_c String CODEC(ZSTD(1)) COMMENT 'example' TTL event_time + INTERVAL 1 DAY,
    value_d String CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY COMMENT 'example',
    value_e String TTL event_time + INTERVAL 1 DAY COMMENT 'example' CODEC(ZSTD(1)),
    value_f String TTL event_time + INTERVAL 1 DAY CODEC(ZSTD(1)) COMMENT 'example',
    value_g String COMMENT 'example' CODEC(ZSTD(1)),
    value_h String CODEC(ZSTD(1)) COMMENT 'example'
)
ENGINE = MergeTree
ORDER BY event_time;

ALTER TABLE db.example ADD COLUMN extra String DEFAULT 'value' TTL event_time + INTERVAL 1 DAY CODEC(LZ4) COMMENT 'example';
ALTER TABLE db.example MODIFY COLUMN extra String MATERIALIZED toString(event_time) CODEC(LZ4) COMMENT 'example' TTL event_time + INTERVAL 1 DAY;


-- Format SQL:
CREATE TABLE db.example (event_time DateTime, value_a String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY, value_b String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY, value_c String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY, value_d String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY, value_e String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY, value_f String COMMENT 'example' CODEC(ZSTD(1)) TTL event_time + INTERVAL 1 DAY, value_g String COMMENT 'example' CODEC(ZSTD(1)), value_h String COMMENT 'example' CODEC(ZSTD(1))) ENGINE = MergeTree ORDER BY event_time;
ALTER TABLE db.example ADD COLUMN extra String DEFAULT 'value' COMMENT 'example' CODEC(LZ4) TTL event_time + INTERVAL 1 DAY;
ALTER TABLE db.example MODIFY COLUMN extra String MATERIALIZED toString(event_time) COMMENT 'example' CODEC(LZ4) TTL event_time + INTERVAL 1 DAY;
