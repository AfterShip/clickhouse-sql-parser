-- Origin SQL:
SELECT number as key, 5 * number value, 'original' AS source
FROM numbers(16)
WHERE (number % 5) == 0
ORDER BY key WITH FILL STALENESS 11;


-- Beautify SQL:
SELECT
  number AS key,
  5 * number AS value,
  'original' AS source
FROM
  numbers(16)
WHERE
  (number % 5) == 0
ORDER BY
  key WITH FILL STALENESS 11;
