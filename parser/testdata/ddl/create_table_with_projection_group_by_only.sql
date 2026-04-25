CREATE TABLE events
(
    `event_time` DateTime,
    `event_type` String,
    `user_id` UInt64,
    `value` Float64,
    PROJECTION hourly_aggregates
    (
        SELECT
            toStartOfHour(event_time) AS hour,
            event_type,
            count() AS event_count,
            sum(value) AS total_value
        GROUP BY hour, event_type
    )
)
ENGINE = MergeTree()
ORDER BY (event_time, event_type);
