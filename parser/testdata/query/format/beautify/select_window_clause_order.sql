-- Origin SQL:
SELECT a, sum(b) OVER w AS s
FROM t
PREWHERE a > 0
WHERE b > 0
GROUP BY a, b
HAVING count() > 1
WINDOW w AS (PARTITION BY a ORDER BY b)
QUALIFY s > 1
ORDER BY a
LIMIT 5;


-- Beautify SQL:
SELECT
  a,
  sum(b) OVER w AS s
FROM
  t
PREWHERE a > 0
WHERE
  b > 0
GROUP BY
  a, b
HAVING count() > 1
WINDOW w AS (PARTITION BY a ORDER BY
  b)
QUALIFY s > 1
ORDER BY
  a
LIMIT 5;
