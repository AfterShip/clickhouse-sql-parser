-- Origin SQL:
SELECT max(end), CASE WHEN end > start THEN 1 ELSE 0 END AS active
FROM t
WHERE end > start
GROUP BY end
ORDER BY end


-- Format SQL:
SELECT max(end), CASE WHEN end > start THEN 1 ELSE 0 END AS active FROM t WHERE end > start GROUP BY end ORDER BY end;
