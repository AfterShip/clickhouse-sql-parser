-- Origin SQL:
SELECT * FROM t1 GLOBAL JOIN t2 ON t1.a = t2.a;
SELECT * FROM t1 GLOBAL INNER JOIN t2 ON t1.a = t2.a;
SELECT * FROM t1 GLOBAL LEFT JOIN t2 ON t1.a = t2.a;
SELECT * FROM t1 GLOBAL LEFT OUTER JOIN t2 USING (a);
SELECT * FROM t1 GLOBAL ANY LEFT JOIN t2 ON t1.a = t2.a;
SELECT * FROM t1 GLOBAL CROSS JOIN t2;
SELECT * FROM t1 AS x LOCAL FULL JOIN t2 ON x.a = t2.a;
SELECT * FROM t1 AS x LOCAL RIGHT JOIN t2 USING a;
SELECT * FROM t1 GLOBAL LEFT JOIN t2 ON t1.a = t2.a GLOBAL LEFT JOIN t3 ON t1.a = t3.a;
SELECT * FROM t WHERE a GLOBAL IN (SELECT b FROM t2);
SELECT * FROM t WHERE a GLOBAL NOT IN (SELECT b FROM t2);
SELECT * FROM t1 GLOBAL LEFT JOIN t2 ON t1.a = t2.a WHERE t1.a GLOBAL NOT IN (SELECT b FROM t3);
SELECT * FROM numbers(3) AS a GLOBAL JOIN numbers(3) AS b ON a.number = b.number;
SELECT a.number FROM numbers(3) AS a GLOBAL LEFT JOIN numbers(2) AS b USING (number);
SELECT * FROM numbers(3) AS a LOCAL ANY LEFT JOIN numbers(3) AS b ON a.number = b.number;
SELECT number FROM numbers(5) WHERE number GLOBAL NOT IN (SELECT number FROM numbers(2));
SELECT * FROM numbers(3) AS global GLOBAL JOIN numbers(3) AS b ON global.number = b.number;


-- Beautify SQL:
SELECT
  *
FROM
  t1
  GLOBAL JOIN
    t2 ON t1.a = t2.a;
SELECT
  *
FROM
  t1
  GLOBAL INNER JOIN
    t2 ON t1.a = t2.a;
SELECT
  *
FROM
  t1
  GLOBAL LEFT JOIN
    t2 ON t1.a = t2.a;
SELECT
  *
FROM
  t1
  GLOBAL LEFT OUTER JOIN
    t2 USING a;
SELECT
  *
FROM
  t1
  GLOBAL ANY LEFT JOIN
    t2 ON t1.a = t2.a;
SELECT
  *
FROM
  t1
  GLOBAL CROSS JOIN
    t2;
SELECT
  *
FROM
  t1 AS x
  LOCAL FULL JOIN
    t2 ON x.a = t2.a;
SELECT
  *
FROM
  t1 AS x
  LOCAL RIGHT JOIN
    t2 USING a;
SELECT
  *
FROM
  t1
  GLOBAL LEFT JOIN
    t2 ON t1.a = t2.a
  GLOBAL LEFT JOIN
    t3 ON t1.a = t3.a;
SELECT
  *
FROM
  t
WHERE
  a GLOBAL IN (SELECT
    b
  FROM
    t2);
SELECT
  *
FROM
  t
WHERE
  a GLOBAL NOT IN (SELECT
    b
  FROM
    t2);
SELECT
  *
FROM
  t1
  GLOBAL LEFT JOIN
    t2 ON t1.a = t2.a
WHERE
  t1.a GLOBAL NOT IN (SELECT
    b
  FROM
    t3);
SELECT
  *
FROM
  numbers(3) AS a
  GLOBAL JOIN
    numbers(3) AS b ON a.number = b.number;
SELECT
  a.number
FROM
  numbers(3) AS a
  GLOBAL LEFT JOIN
    numbers(2) AS b USING number;
SELECT
  *
FROM
  numbers(3) AS a
  LOCAL ANY LEFT JOIN
    numbers(3) AS b ON a.number = b.number;
SELECT
  number
FROM
  numbers(5)
WHERE
  number GLOBAL NOT IN (SELECT
    number
  FROM
    numbers(2));
SELECT
  *
FROM
  numbers(3) AS global
  GLOBAL JOIN
    numbers(3) AS b ON global.number = b.number;
