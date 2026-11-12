SELECT id, row_number() OVER (PARTITION BY id ORDER BY ts) AS rn FROM t QUALIFY rn = 1;
SELECT id
FROM t
WHERE id > 0
GROUP BY id
HAVING count() > 1
QUALIFY row_number() OVER (ORDER BY id) = 1
ORDER BY id
LIMIT 5;
SELECT id FROM t WINDOW w AS (PARTITION BY id) QUALIFY row_number() OVER w = 1;
SELECT qualify FROM t AS qualify QUALIFY qualify = 1
