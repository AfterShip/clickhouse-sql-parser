-- Origin SQL:
SELECT t.a.1, db.t.a.2 + 1 FROM db.t AS t WHERE t.a.1 = 1 ORDER BY t.a.2


-- Beautify SQL:
SELECT
  t.a.1,
  db.t.a.2 + 1
FROM
  db.t AS t
WHERE
  t.a.1 = 1
ORDER BY
  t.a.2;
