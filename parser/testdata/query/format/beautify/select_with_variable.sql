-- Origin SQL:
WITH $abc AS (SELECT 1 AS a) SELECT * FROM $abc

-- Beautify SQL:
WITH
  $abc AS (SELECT
    1 AS a)
SELECT
  *
FROM
  $abc;
