-- Origin SQL:
SELECT
    v,
    j
FROM t1
    ARRAY JOIN JSONExtractArrayRaw(a) AS j
    ARRAY JOIN array(
    JSONExtractString(j, 'x'),
    JSONExtractString(j, 'y')
) AS v;

-- Beautify SQL:
SELECT
  v,
  j
FROM
  t1
  ARRAY JOIN
    JSONExtractArrayRaw(a) AS j
  ARRAY JOIN
    array(JSONExtractString(j, 'x'), JSONExtractString(j, 'y')) AS v;
