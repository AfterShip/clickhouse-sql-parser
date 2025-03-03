-- Origin SQL:
SELECT aggregation_target AS aggregation_target,
    timestamp AS timestamp,
    step_0 AS step_0,
    latest_0 AS latest_0,
    step_1 AS step_1,
    latest_1 AS latest_1,
    step_2 AS step_2,
    min(latest_2) OVER (PARTITION BY aggregation_target
    ORDER BY timestamp DESC ROWS BETWEEN UNBOUNDED PRECEDING AND 0 PRECEDING) AS latest_2
FROM t0

-- Format SQL:
SELECT aggregation_target AS aggregation_target, timestamp AS timestamp, step_0 AS step_0, latest_0 AS latest_0, step_1 AS step_1, latest_1 AS latest_1, step_2 AS step_2, min(latest_2) OVER ( PARTITION BY aggregation_target ORDER BY timestamp DESC ROWS  BETWEEN  UNBOUNDED PRECEDING AND  0 PRECEDING) AS latest_2 FROM t0;
