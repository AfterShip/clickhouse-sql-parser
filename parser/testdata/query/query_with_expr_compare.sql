SELECT date, path, splitByChar('/', path)[2] AS path_b
FROM(
    SELECT 'pathA/pathB/pathC' AS path, '2024-09-10' AS date
    )
WHERE toDate(date) BETWEEN '2024-09-01' AND '2024-09-30'
  AND splitByChar('/', path)[1] = 'pathA'