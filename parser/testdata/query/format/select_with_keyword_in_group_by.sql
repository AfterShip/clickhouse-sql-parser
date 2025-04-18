-- Origin SQL:
SELECT
    toStartOfInterval(timestamp, toIntervalMinute(1)) AS interval,
    column_name
FROM table
WHERE true
GROUP BY (interval, column_name)
ORDER BY (interval AS i, column_name) ASC

-- Format SQL:
SELECT toStartOfInterval(timestamp, toIntervalMinute(1)) AS interval, column_name FROM table WHERE true GROUP BY (interval, column_name) ORDER BY (interval AS i, column_name) ASC;
