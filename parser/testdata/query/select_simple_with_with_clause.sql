WITH
    cte1 AS (SELECT f1 FROM t1),
    cte2 AS (SELECT f2 FROM t2)
SELECT
    cte1.f1,
    cte2.f2,
    t3.f3
FROM
    t3,cte1,cte2

