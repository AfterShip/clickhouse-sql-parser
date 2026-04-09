-- Origin SQL:
CREATE MATERIALIZED VIEW db.mv_with_comment TO db.dst_table
(
    `shop_id` UInt64,
    `event_type` LowCardinality(String),
    `created_at` DateTime64(9)
)
COMMENT '{"blueprint_hash":"abc123","timestamp":"2026-04-08T12:00:00Z"}'
AS SELECT
    shop_id,
    event_type,
    created_at
FROM db.src_table;


-- Format SQL:
CREATE MATERIALIZED VIEW db.mv_with_comment TO db.dst_table (`shop_id` UInt64, `event_type` LowCardinality(String), `created_at` DateTime64(9)) AS SELECT shop_id, event_type, created_at FROM db.src_table COMMENT '{"blueprint_hash":"abc123","timestamp":"2026-04-08T12:00:00Z"}';
