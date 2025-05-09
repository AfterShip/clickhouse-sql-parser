-- Origin SQL:
CREATE MATERIALIZED VIEW fresh_mv
REFRESH EVERY 1 HOUR OFFSET 10 MINUTE APPEND TO events_export
(
    `timestamp` DateTime64(9),
    `field_1` String,
    `field_2` String,
)
DEFINER = default SQL SECURITY DEFINER
AS (SELECT
    timestamp,
    field_1,
    field_2,
FROM event_table
WHERE toStartOfHour(timestamp) = toStartOfHour(now() - toIntervalHour(1)))
COMMENT 'Test comment'


-- Format SQL:
CREATE MATERIALIZED VIEW fresh_mv REFRESH EVERY 1 HOUR OFFSET 10 MINUTE APPEND TO events_export (`timestamp` DateTime64(9), `field_1` String, `field_2` String) DEFINER = default SQL SECURITY DEFINER AS (SELECT timestamp, field_1, field_2, FROM AS event_table WHERE toStartOfHour(timestamp) = toStartOfHour(now() - toIntervalHour(1))) COMMENT 'Test comment';
