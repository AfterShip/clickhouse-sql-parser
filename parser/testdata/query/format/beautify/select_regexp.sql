-- Origin SQL:
SELECT a, b
FROM t
WHERE name REGEXP '^foo'


-- Beautify SQL:
SELECT
  a,
  b
FROM
  t
WHERE
  name REGEXP '^foo';
