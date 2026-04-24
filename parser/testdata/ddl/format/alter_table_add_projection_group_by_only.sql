-- Origin SQL:
ALTER TABLE events
ADD PROJECTION IF NOT EXISTS hourly_stats
(SELECT toStartOfHour(event_time) AS hour, event_type, count() AS count, uniq(user_id) AS users GROUP BY hour, event_type);


-- Format SQL:
ALTER TABLE events ADD PROJECTION IF NOT EXISTS hourly_stats (SELECT toStartOfHour(event_time) AS hour, event_type, count() AS count, uniq(user_id) AS users GROUP BY hour, event_type);
