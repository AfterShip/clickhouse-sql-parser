SELECT
    f0, coalesce(f1, f2) AS f3, row_number()
OVER (PARTITION BY f0 ORDER BY f1 ASC) AS rn
FROM test.events_local
WHERE (f0 IN ('foo', 'bar', 'test')) AND (f1 = 'testing') AND (f2 NOT LIKE 'testing2')
AND f3 NOT IN ('a', 'b', 'c')


GROUP BY f0,   f1

Limit 100, 10 By f0;