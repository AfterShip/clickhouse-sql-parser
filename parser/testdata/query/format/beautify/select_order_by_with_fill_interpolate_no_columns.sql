-- Origin SQL:
SELECT n, value FROM (
   SELECT toFloat32(number % 10) AS n, number AS value
   FROM numbers(10) WHERE number % 3 = 1
) ORDER BY n WITH FILL FROM 0 TO 10 STEP 1
INTERPOLATE;


-- Beautify SQL:
SELECT
  n,
  value
FROM
  (SELECT
    toFloat32(number % 10) AS n,
    number AS value
  FROM
    numbers(10)
  WHERE
    number % 3 = 1)
ORDER BY
  n WITH FILL FROM 0 TO 10 STEP 1
  INTERPOLATE;
