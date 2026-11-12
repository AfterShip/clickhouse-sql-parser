-- Origin SQL:
SELECT a.id
FROM t1 AS a
ANTI LEFT JOIN (SELECT id FROM t2) AS b ON a.id = b.id
ANTI RIGHT JOIN t3 AS c ON a.id = c.id
SEMI LEFT JOIN t4 AS d ON a.id = d.id
LEFT ANTI JOIN t5 AS e ON a.id = e.id
GLOBAL ANTI LEFT JOIN t6 AS f ON a.id = f.id


-- Beautify SQL:
SELECT
  a.id
FROM
  t1 AS a
  ANTI LEFT JOIN
    (SELECT
      id
    FROM
      t2) AS b ON a.id = b.id
  ANTI RIGHT JOIN
    t3 AS c ON a.id = c.id
  SEMI LEFT JOIN
    t4 AS d ON a.id = d.id
  LEFT ANTI JOIN
    t5 AS e ON a.id = e.id
  GLOBAL ANTI LEFT JOIN
    t6 AS f ON a.id = f.id;
