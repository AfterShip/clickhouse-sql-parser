-- Origin SQL:
SELECT n, source, inter FROM (
   SELECT toFloat32(number % 10) AS n, 'original' AS source, number AS inter
   FROM numbers(10) WHERE number % 3 = 1
) ORDER BY n WITH FILL FROM 0 TO 5.51 STEP 0.5
INTERPOLATE (inter AS inter + 1);


-- Beautify SQL:
SELECT
  n,
  source,
  inter
FROM
  (SELECT
    toFloat32(number % 10) AS n,
    'original' AS source,
    number AS inter
  FROM
    numbers(10)
  WHERE
    number % 3 = 1)
ORDER BY
  n WITH FILL FROM 0 TO 5.51 STEP 0.5
  INTERPOLATE (inter AS inter + 1);
