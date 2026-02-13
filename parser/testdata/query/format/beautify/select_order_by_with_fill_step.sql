-- Origin SQL:
SELECT date, value FROM (
    SELECT toDate('2020-01-01') + INTERVAL number DAY AS date, number AS value
    FROM numbers(5)
) ORDER BY date WITH FILL STEP INTERVAL 1 DAY;


-- Beautify SQL:
SELECT
  date,
  value
FROM
  (SELECT
    toDate('2020-01-01') + INTERVAL number DAY AS date,
    number AS value
  FROM
    numbers(5))
ORDER BY
  date WITH FILL STEP INTERVAL 1 DAY;
