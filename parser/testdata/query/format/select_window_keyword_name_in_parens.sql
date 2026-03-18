-- Origin SQL:
SELECT sum(x) OVER (order) AS sum_over_order
FROM t
WINDOW order AS (PARTITION BY team ORDER BY ts);


-- Format SQL:
SELECT sum(x) OVER (order) AS sum_over_order FROM t WINDOW order AS (PARTITION BY team ORDER BY ts);
