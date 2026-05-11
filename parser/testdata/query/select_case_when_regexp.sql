SELECT
    CASE WHEN col REGEXP '^[0-9]+$' THEN toInt32(col) ELSE 0 END AS num_value
FROM t
