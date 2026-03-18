SELECT sum(x) OVER (w1 ORDER BY ts ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) AS rolling_sum,
       avg(x) OVER (w2)                                                      AS avg_over_w2
FROM t
WINDOW w1 AS (PARTITION BY team),
       w2 AS (w1 ORDER BY ts ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW);
