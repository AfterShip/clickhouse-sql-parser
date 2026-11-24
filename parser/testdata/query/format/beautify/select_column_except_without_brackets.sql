-- Origin SQL:
SELECT * EXCEPT b FROM t;
SELECT t.* EXCEPT b, c FROM t;
SELECT * EXCEPT b EXCEPT (c, d) FROM t;


-- Beautify SQL:
SELECT
  * EXCEPT(b)
FROM
  t;
SELECT
  t.* EXCEPT(b),
  c
FROM
  t;
SELECT
  * EXCEPT(b) EXCEPT(c, d)
FROM
  t;
