-- Origin SQL:
SELECT a.id
FROM t1 AS a
ANTI LEFT JOIN t2 AS b ON a.id = b.id
ANTI RIGHT JOIN t3 AS c ON a.id = c.id
GLOBAL ANTI LEFT OUTER JOIN t4 AS d ON a.id = d.id
ANTI JOIN t5 AS e ON a.id = e.id


-- Format SQL:
SELECT a.id FROM t1 AS a ANTI LEFT JOIN t2 AS b ON a.id = b.id ANTI RIGHT JOIN t3 AS c ON a.id = c.id GLOBAL ANTI LEFT OUTER JOIN t4 AS d ON a.id = d.id ANTI JOIN t5 AS e ON a.id = e.id;
