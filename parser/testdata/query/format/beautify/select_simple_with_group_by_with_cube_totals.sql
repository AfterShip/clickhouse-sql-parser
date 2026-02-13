-- Origin SQL:
SELECT a, COUNT(b) FROM group_by_all GROUP BY CUBE(a) WITH CUBE WITH TOTALS ORDER BY a;

-- Beautify SQL:
SELECT
  a,
  COUNT(b)
FROM
  group_by_all
GROUP BY
  CUBE(a)
  WITH CUBE
  WITH TOTALS
ORDER BY
  a;
