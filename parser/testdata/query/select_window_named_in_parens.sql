SELECT sum(x) OVER (w) AS sum_over_w
FROM t
WINDOW w AS (PARTITION BY y ORDER BY x);
